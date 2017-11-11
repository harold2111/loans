package loan

import (
	"loans/config"
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

func (bill *Bill) applyPayment(paymentToBill decimal.Decimal) {
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
