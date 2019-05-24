package models

import (
	"loans/shared/config"
	"loans/shared/utils"
	"loans/shared/utils/financial"
	"time"

	"github.com/shopspring/decimal"
)

const (
	BillStateDue       = "DUE"
	BillStatePaid      = "PAID"
	PeriodStatusClosed = "ClOSED"
	PeriodStatusOpen   = "OPEN"
)

type Bill struct {
	ID                  uint            `gorm:"primary_key" json:"id"`
	CreatedAt           time.Time       `json:"-"`
	UpdatedAt           time.Time       `json:"-"`
	DeletedAt           *time.Time      `sql:"index" json:"-"`
	LoanID              uint            `json:"loanID"`
	State               string          `json:"state"`
	PeriodStatus        string          `json:"periodStatus"`
	Period              uint            `json:"period"`
	BillStartDate       time.Time       `json:"billStartDate"`
	BillEndDate         time.Time       `json:"billEndDate"`
	PaymentDate         time.Time       `json:"paymentDate"`
	InitialPrincipal    decimal.Decimal `gorm:"type:numeric" json:"initialPrincipal"`
	Payment             decimal.Decimal `gorm:"type:numeric" json:"payment"`
	InterestRate        decimal.Decimal `gorm:"type:numeric" json:"interestRate"`
	InterestOfPayment   decimal.Decimal `gorm:"type:numeric" json:"interestOfPayment"`
	PrincipalOfPayment  decimal.Decimal `gorm:"type:numeric" json:"principalOfPayment"`
	Paid                decimal.Decimal `gorm:"type:numeric" json:"paid"`
	DaysLate            int             `json:"daysLate"`
	FeeLateDue          decimal.Decimal `gorm:"type:numeric" json:"feeLateDue"`
	PaymentDue          decimal.Decimal `gorm:"type:numeric" json:"paymentDue"`
	TotalDue            decimal.Decimal `gorm:"type:numeric" json:"totalDue"`
	PaidToPrincipal     decimal.Decimal `gorm:"type:numeric" json:"paidToPrincipal"`
	FinalPrincipal      decimal.Decimal `gorm:"type:numeric" json:"finalPrincipal"`
	LastLiquidationDate time.Time       `json:"lastLiquidationDate"`
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
