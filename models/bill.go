package models

import (
	"fmt"
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
	BillStartDate       time.Time       `gorm:"type:timestamp without time zone"`
	BillEndDate         time.Time       `gorm:"type:timestamp without time zone"`
	PaymentDate         time.Time       `gorm:"type:timestamp without time zone"`
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
	LastLiquidationDate time.Time       `gorm:"type:timestamp without time zone"`
}

func (loanBill *Bill) Create() error {
	error := config.DB.Create(loanBill).Error
	return error
}

func (loanBill *Bill) Update() error {
	error := config.DB.Save(loanBill).Error
	return error
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

func FindBillOpenPeriodByLoanID(loanID uint) (*Bill, error) {
	bill := Bill{}
	error := config.DB.Raw("SELECT * FROM bills WHERE loan_id = ? AND period_status = ? AND period = (SELECT max(period) FROM bills where loan_id = ?)",
		loanID, PeriodStatusOpen, loanID).Scan(&bill).Error
	return &bill, error
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
	round := config.Round
	balance := balanceExpectedInSpecificPeriodOfLoan(loan, period)
	newBill := Bill{}
	newBill.LoanID = loanID
	newBill.State = BillStateDue
	newBill.PeriodStatus = PeriodStatusOpen
	newBill.Period = uint(period)
	newBill.BillStartDate = loan.StartDate
	newBill.BillEndDate = utils.AddMothToTimeUtil(newBill.BillStartDate, 1)
	newBill.PaymentDate = newBill.BillEndDate
	newBill.Payment = balance.Payment.RoundBank(round)
	newBill.InterestOfPayment = balance.ToInterest.RoundBank(round)
	newBill.InterestRate = loan.InterestRatePeriod.RoundBank(round)
	newBill.PrincipalOfPayment = balance.ToPrincipal.RoundBank(round)
	newBill.Paid = decimal.Zero
	newBill.DaysLate = 0
	newBill.FeeLateDue = decimal.Zero
	newBill.PaymentDue = newBill.Payment
	newBill.TotalDue = newBill.Payment
	newBill.LastLiquidationDate = newBill.PaymentDate
	newBill.Create()

	return nil

}

func RecurringLoanBillingByLoanID(loanID uint) error {
	loan, error := FindLoanByID(loanID)
	if error != nil {
		return error
	}
	oldLoanBill := new(Bill)
	oldLoanBill, error = FindBillOpenPeriodByLoanID(loanID)
	if error != nil {
		return error
	}
	now := time.Now().In(oldLoanBill.BillEndDate.Location())
	if now.Before(oldLoanBill.BillEndDate) {
		return nil
	}
	period := int(oldLoanBill.Period + 1)
	round := config.Round
	balance := balanceExpectedInSpecificPeriodOfLoan(loan, period)
	newBill := Bill{}
	newBill.LoanID = loanID
	newBill.State = BillStateDue
	newBill.PeriodStatus = PeriodStatusOpen
	newBill.Period = uint(period)
	newBill.BillStartDate = oldLoanBill.BillEndDate.AddDate(0, 0, 1)
	newBill.BillEndDate = utils.AddMothToTimeUtil(newBill.BillStartDate, 1)
	newBill.PaymentDate = newBill.BillEndDate
	newBill.Payment = balance.Payment.RoundBank(round)
	newBill.InterestOfPayment = balance.ToInterest.RoundBank(round)
	newBill.InterestRate = loan.InterestRatePeriod.RoundBank(round)
	newBill.PrincipalOfPayment = balance.ToPrincipal.RoundBank(round)
	newBill.Paid = decimal.Zero
	newBill.DaysLate = 0
	newBill.FeeLateDue = decimal.Zero
	newBill.PaymentDue = newBill.Payment
	newBill.TotalDue = newBill.Payment
	newBill.LastLiquidationDate = newBill.PaymentDate
	newBill.Create()
	oldLoanBill.PeriodStatus = PeriodStatusClosed
	oldLoanBill.Update()
	return nil
}

func (bill *Bill) LiquidateBill() {
	fmt.Println("time.Now():", time.Now())
	location, _ := time.LoadLocation(config.DefaultLocation)
	fmt.Println("bill.LastLiquidationDate:", bill.LastLiquidationDate)
	fmt.Println("bill.LastLiquidationDate.IN:", bill.LastLiquidationDate.In(location))
	now := time.Now()
	daysLate := 0
	if now.After(bill.PaymentDate) {
		daysLate = utils.DaysSince(bill.LastLiquidationDate)
		if daysLate < 0 {
			daysLate = 0
		}
	}
	feeLatePeriod := financial.FeeLateWithPeriodInterest(bill.InterestRate, bill.PaymentDue, daysLate).RoundBank(config.Round)
	totalFeeLateDue := bill.FeeLateDue.Add(feeLatePeriod).RoundBank(config.Round)
	totalDue := bill.PaymentDue.Add(totalFeeLateDue).RoundBank(config.Round)
	totalDaysLate := bill.DaysLate + daysLate

	bill.DaysLate = totalDaysLate
	bill.FeeLateDue = totalFeeLateDue
	bill.TotalDue = totalDue
	bill.LastLiquidationDate = now
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
}
