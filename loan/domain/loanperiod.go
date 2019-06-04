package domain

import (
	"loans/shared/config"
	"loans/shared/utils"
	"loans/shared/utils/financial"
	"time"

	"github.com/shopspring/decimal"
)

const (
	LoanPeriodStateDue    = "DUE"
	LoanPeriodStatePaid   = "PAID"
	LoanPeriodStateClosed = "ClOSED"
	PeriodStatusOpen      = "OPEN"
)

type LoanPeriod struct {
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

func (period *LoanPeriod) LiquidateBill(liquidationDate time.Time) {
	daysLate := calculateDaysLate(period.LastLiquidationDate, liquidationDate)
	feeLatePeriod := financial.FeeLateWithPeriodInterest(period.InterestRate, period.PaymentDue, daysLate).RoundBank(config.Round)
	totalFeeLateDue := period.FeeLateDue.Add(feeLatePeriod).RoundBank(config.Round)
	totalDue := period.PaymentDue.Add(totalFeeLateDue).RoundBank(config.Round)
	totalDaysLate := period.DaysLate + daysLate

	period.DaysLate = totalDaysLate
	period.FeeLateDue = totalFeeLateDue
	period.TotalDue = totalDue
	period.LastLiquidationDate = liquidationDate
}

func (period *LoanPeriod) ApplyPayment(paymentToBill decimal.Decimal) {
	//the payment NO covers all the fee late
	if paymentToBill.LessThanOrEqual(period.FeeLateDue) {
		period.FeeLateDue = period.FeeLateDue.Sub(paymentToBill)
	} else { //the payment covers fee late
		remainingPaymentToBill := paymentToBill.Sub(period.FeeLateDue)
		period.FeeLateDue = decimal.Zero
		paymentDue := period.PaymentDue.Sub(remainingPaymentToBill).RoundBank(config.Round)
		if paymentDue.LessThanOrEqual(decimal.Zero) {
			period.PaidToPrincipal = period.PaidToPrincipal.Add(paymentDue.Abs()).RoundBank(config.Round)
			period.FinalPrincipal = period.FinalPrincipal.Sub(period.PaidToPrincipal).RoundBank(config.Round)
			period.PaymentDue = decimal.Zero
		} else {
			period.PaymentDue = paymentDue
		}
	}
	period.TotalDue = period.PaymentDue.Add(period.FeeLateDue).RoundBank(config.Round)
	period.Paid = period.Paid.Add(paymentToBill)
	if period.TotalDue.LessThanOrEqual(decimal.Zero) {
		period.State = LoanPeriodStatePaid
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
