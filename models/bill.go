package models

import (
	"loans/config"
	"loans/errors"
	"loans/financial"
	"loans/utils"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

const (
	BillStateDue       = "DUE"
	BillStatePaid      = "PAID"
	PeriodStatusClosed = "ClOSED"
	PeriodStatusOpen   = "OPEN"
)

type Bill struct {
	gorm.Model
	LoanID              uint
	State               string
	PeriodStatus        string
	Period              uint
	BillStartDate       time.Time
	BillEndDate         time.Time
	PaymentDate         time.Time
	InitialPrincipal    decimal.Decimal `gorm:"type:numeric"`
	Payment             decimal.Decimal `gorm:"type:numeric"`
	InterestRate        decimal.Decimal `gorm:"type:numeric"`
	InterestOfPayment   decimal.Decimal `gorm:"type:numeric"`
	PrincipalOfPayment  decimal.Decimal `gorm:"type:numeric"`
	Paid                decimal.Decimal `gorm:"type:numeric"`
	DaysLate            int
	FeeLateDue          decimal.Decimal `gorm:"type:numeric"`
	PaymentDue          decimal.Decimal `gorm:"type:numeric"`
	TotalDue            decimal.Decimal `gorm:"type:numeric"`
	PaidToPrincipal     decimal.Decimal `gorm:"type:numeric"`
	FinalPrincipal      decimal.Decimal `gorm:"type:numeric"`
	LastLiquidationDate time.Time
}

func (bill *Bill) Create() error {
	error := config.DB.Create(bill).Error
	return error
}

func (bill *Bill) Update() error {
	error := config.DB.Save(bill).Error
	return error
}

func (bill *Bill) ClosePeriod() error {
	bill.PeriodStatus = PeriodStatusClosed
	return bill.Update()
}

func FindBillsByLoanID(loanID uint) ([]Bill, error) {
	var bills []Bill
	config.DB.Find(&bills, "loan_id = ?", loanID)
	return bills, nil
}

func FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID uint) ([]Bill, error) {
	var bills []Bill
	config.DB.Order("period").Find(&bills, "loan_id = ? AND state = ? OR period_status = ?", loanID, BillStateDue, PeriodStatusOpen)
	return bills, nil
}

func FindBillOpenPeriodByLoanID(loanID uint) (Bill, error) {
	bill := Bill{}
	error := config.DB.Raw("SELECT * FROM bills WHERE loan_id = ? AND period_status = ? AND period = (SELECT max(period) FROM bills where loan_id = ?)",
		loanID, PeriodStatusOpen, loanID).Scan(&bill).Error
	return bill, error
}

func CreateInitialBill(loanID uint) error {
	loan, error := FindLoanByID(loanID)
	if error != nil {
		return error
	}
	bills, _ := FindBillsByLoanID(loanID)
	if len(bills) > 0 {
		return &errors.GracefulError{ErrorCode: errors.BillAlreadyExist}
	}
	period := 1
	newBill := Bill{}
	newBill.LoanID = loan.ID
	newBill.Period = uint(period)
	newBill.BillStartDate = loan.StartDate
	newBill.BillEndDate = utils.AddMothToTimeForPayment(newBill.BillStartDate, 1)
	balancePeriod := balanceExpectedInSpecificPeriodOfLoan(loan, period)
	fillDefaultAmountValues(&newBill, balancePeriod)
	newBill.Create()

	return nil
}

func RecurringLoanBillingByLoanID(loanID uint) error {
	loan, error := FindLoanByID(loanID)
	if error != nil {
		return error
	}
	oldBill := Bill{}
	oldBill, error = FindBillOpenPeriodByLoanID(loanID)
	if error != nil {
		return error
	}
	if time.Now().Before(oldBill.BillEndDate) {
		return nil
	}
	period := int(oldBill.Period + 1)
	newBill := Bill{}
	newBill.LoanID = loan.ID
	newBill.Period = uint(period)
	newBill.BillStartDate = oldBill.BillEndDate.AddDate(0, 0, 1)
	newBill.BillEndDate = utils.AddMothToTimeForPayment(oldBill.BillEndDate, 1)
	nextBalance := nextBalanceFromBill(oldBill)
	fillDefaultAmountValues(&newBill, nextBalance)
	if error := newBill.Create(); error != nil {
		return error
	}
	if error := oldBill.ClosePeriod(); error != nil {
		return error
	}
	return RecurringLoanBillingByLoanID(loanID)
}

func nextBalanceFromBill(bill Bill) financial.Balance {
	balance := financial.Balance{}
	balance.InitialPrincipal = bill.InitialPrincipal
	balance.Payment = bill.Payment
	balance.InterestRatePeriod = bill.InterestRate
	balance.ToInterest = bill.InterestOfPayment
	balance.ToPrincipal = bill.PrincipalOfPayment
	balance.FinalPrincipal = bill.FinalPrincipal
	return financial.NextBalanceFromBefore(balance)
}

func fillDefaultAmountValues(bill *Bill, balance financial.Balance) {
	round := config.Round
	bill.State = BillStateDue
	bill.PeriodStatus = PeriodStatusOpen
	bill.PaymentDate = bill.BillEndDate
	bill.InitialPrincipal = balance.InitialPrincipal
	bill.Payment = balance.Payment.RoundBank(round)
	bill.InterestOfPayment = balance.ToInterest.RoundBank(round)
	bill.InterestRate = balance.InterestRatePeriod.RoundBank(round)
	bill.PrincipalOfPayment = balance.ToPrincipal.RoundBank(round)
	bill.Paid = decimal.Zero
	bill.DaysLate = 0
	bill.FeeLateDue = decimal.Zero
	bill.PaymentDue = bill.Payment
	bill.TotalDue = bill.Payment
	bill.FinalPrincipal = balance.FinalPrincipal.RoundBank(round)
	bill.LastLiquidationDate = bill.PaymentDate
}

func (bill *Bill) LiquidateBill(liquidationDate time.Time) {
	daysLate := calculateDaysLate(bill.LastLiquidationDate, liquidationDate)
	feeLatePeriod := financial.FeeLateWithPeriodInterest(bill.InterestRate, bill.PaymentDue, daysLate).RoundBank(config.Round)
	totalFeeLateDue := bill.FeeLateDue.Add(feeLatePeriod).RoundBank(config.Round)
	totalDue := bill.PaymentDue.Add(totalFeeLateDue).RoundBank(config.Round)
	totalDaysLate := bill.DaysLate + daysLate

	bill.DaysLate = totalDaysLate
	bill.FeeLateDue = totalFeeLateDue
	bill.TotalDue = totalDue
	bill.LastLiquidationDate = liquidationDate
}

func (bill *Bill) ApplyPayment(paymentToBill decimal.Decimal) {
	//the payment NO covers all the fee late
	if paymentToBill.LessThanOrEqual(bill.FeeLateDue) {
		bill.FeeLateDue = bill.FeeLateDue.Sub(paymentToBill)
	} else { //the payment covers fee late
		remainingPaymentToBill := paymentToBill.Sub(bill.FeeLateDue)
		bill.FeeLateDue = decimal.Zero
		paymentDue := bill.PaymentDue.Sub(remainingPaymentToBill).RoundBank(config.Round)
		if paymentDue.LessThanOrEqual(decimal.Zero) {
			bill.PaidToPrincipal = bill.PaidToPrincipal.Add(paymentDue.Abs()).RoundBank(config.Round)
			bill.FinalPrincipal = bill.FinalPrincipal.Sub(bill.PaidToPrincipal).RoundBank(config.Round)
			bill.PaymentDue = decimal.Zero
		} else {
			bill.PaymentDue = paymentDue
		}
	}
	bill.TotalDue = bill.PaymentDue.Add(bill.FeeLateDue).RoundBank(config.Round)
	bill.Paid = bill.Paid.Add(paymentToBill)
	if bill.TotalDue.LessThanOrEqual(decimal.Zero) {
		bill.State = BillStatePaid
	}
	if bill.FinalPrincipal.LessThanOrEqual(decimal.Zero) {
		bill.PeriodStatus = PeriodStatusClosed
		CloseLoan(bill.LoanID)
	}
}

func calculateDaysLate(lastLiquidationDate, liquidationDate time.Time) int {
	daysLate := 0
	if liquidationDate.After(lastLiquidationDate) {
		daysLate = utils.DaysBetween(lastLiquidationDate, liquidationDate)
		if daysLate < 0 {
			daysLate = 0
		}
	}
	return daysLate
}
