package domain

import (
	"github.com/harold2111/loans/shared/config"
	"github.com/harold2111/loans/shared/utils"
	"github.com/harold2111/loans/shared/utils/financial"
	"time"

	"github.com/shopspring/decimal"
)

const (
	LoanPeriodStateOpen      = "OPEN"
	LoanPeriodStateDue       = "DUE"
	LoanPeriodStatePaid      = "PAID"
	LoanPeriodStateAnnuelled = "ANNULLED"
)

type LoanPeriod struct {
	ID                 uint `gorm:"primary_key"`
	LoanID             uint
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time `sql:"index"`
	PeriodNumber       uint
	State              string
	StartDate          time.Time
	EndDate            time.Time
	PaymentDate        time.Time
	InitialPrincipal   decimal.Decimal `gorm:"type:numeric"`
	Payment            decimal.Decimal `gorm:"type:numeric"`
	InterestRate       decimal.Decimal `gorm:"type:numeric"`
	PrincipalOfPayment decimal.Decimal `gorm:"type:numeric"`
	InterestOfPayment  decimal.Decimal `gorm:"type:numeric"`
	FinalPrincipal     decimal.Decimal `gorm:"type:numeric"`
	//Modifible fields
	LastLiquidationDate                time.Time
	DaysInArrearsSinceLastLiquidation  int
	DebtForArrearsSinceLastLiquidation decimal.Decimal `gorm:"type:numeric"`
	TotalDaysInArrears                 int
	TotalDebtForArrears                decimal.Decimal `gorm:"type:numeric"`
	TotalDebtOfPayment                 decimal.Decimal `gorm:"type:numeric"`
	TotalDebt                          decimal.Decimal `gorm:"type:numeric"`
	TotalPaid                          decimal.Decimal `gorm:"type:numeric"`
	TotalPaidToDebtForArrears          decimal.Decimal `gorm:"type:numeric"`
	TotalPaidToRegularDebt             decimal.Decimal `gorm:"type:numeric"`
	TotalPaidExtraToPrincipal          decimal.Decimal `gorm:"type:numeric"`
}

func (period *LoanPeriod) LiquidateByDate(liquidationDate time.Time) {
	if liquidationDate.After(period.PaymentDate) && period.State == LoanPeriodStateDue {
		daysInArrearsSinceLastLiquidation := calculateDaysLate(period.LastLiquidationDate, liquidationDate)
		debtForArrearsSinceLastLiquidation := financial.FeeLateWithPeriodInterest(period.InterestRate, period.TotalDebtOfPayment, daysInArrearsSinceLastLiquidation).RoundBank(config.Round)
		period.LastLiquidationDate = liquidationDate
		period.DaysInArrearsSinceLastLiquidation = daysInArrearsSinceLastLiquidation
		period.DebtForArrearsSinceLastLiquidation = debtForArrearsSinceLastLiquidation
		period.TotalDaysInArrears = period.TotalDaysInArrears + daysInArrearsSinceLastLiquidation
		period.TotalDebtForArrears = period.TotalDebtForArrears.Add(debtForArrearsSinceLastLiquidation).RoundBank(config.Round)
		period.TotalDebt = period.TotalDebtOfPayment.Add(period.TotalDebtForArrears)
	}
}

func (period *LoanPeriod) ApplyPayment(paymentID uint, payment decimal.Decimal) LoanPeriodMovement {
	var periodMovement LoanPeriodMovement
	periodMovement.fillInitialMovementFromPeriod(*period)
	periodMovement.PaymentID = paymentID
	periodMovement.Paid = payment

	remainingPayment := payment
	remainingPayment = period.applyPaymentToDebtForArrears(&periodMovement, remainingPayment)
	remainingPayment = period.applyPaymentToDebtOfPayments(&periodMovement, remainingPayment)
	remainingPayment = period.applyPaymentToExtraPrincipal(&periodMovement, remainingPayment)

	period.TotalDebt = period.TotalDebtForArrears.Add(period.TotalDebtOfPayment).RoundBank(config.Round)
	period.TotalPaid = period.TotalPaid.Add(payment)
	if period.TotalDebt.LessThanOrEqual(decimal.Zero) {
		period.State = LoanPeriodStatePaid
	}
	periodMovement.fillFinalMovementFromPeriod(*period)
	return periodMovement
}

func (period *LoanPeriod) applyPaymentToDebtForArrears(periodMovement *LoanPeriodMovement, payment decimal.Decimal) decimal.Decimal {
	remainingPayment := payment
	if period.TotalDebtForArrears.LessThanOrEqual(decimal.Zero) || remainingPayment.LessThanOrEqual(decimal.Zero) {
		return remainingPayment
	}
	var paymentToDebtForArrears decimal.Decimal
	if remainingPayment.LessThanOrEqual(period.TotalDebtForArrears) {
		paymentToDebtForArrears = remainingPayment
	} else {
		paymentToDebtForArrears = period.TotalDebtForArrears
	}
	periodMovement.PaidToDebtForArrears = paymentToDebtForArrears
	period.TotalPaidToDebtForArrears = period.TotalPaidToDebtForArrears.Add(paymentToDebtForArrears).RoundBank(config.Round)
	period.TotalDebtForArrears = period.TotalDebtForArrears.Sub(paymentToDebtForArrears).RoundBank(config.Round)
	return remainingPayment.Sub(paymentToDebtForArrears).RoundBank(config.Round)
}

func (period *LoanPeriod) applyPaymentToDebtOfPayments(periodMovement *LoanPeriodMovement, payment decimal.Decimal) decimal.Decimal {
	remainingPayment := payment
	if period.TotalDebtOfPayment.LessThanOrEqual(decimal.Zero) || remainingPayment.LessThanOrEqual(decimal.Zero) {
		return remainingPayment
	}
	var paymentToRegularDebt decimal.Decimal
	if remainingPayment.LessThanOrEqual(period.TotalDebtOfPayment) {
		paymentToRegularDebt = remainingPayment
	} else {
		paymentToRegularDebt = period.TotalDebtOfPayment
	}
	periodMovement.PaidToPaymentDebt = paymentToRegularDebt
	period.TotalPaidToRegularDebt = period.TotalPaidToRegularDebt.Add(paymentToRegularDebt).RoundBank(config.Round)
	period.TotalDebtOfPayment = period.TotalDebtOfPayment.Sub(paymentToRegularDebt).RoundBank(config.Round)
	return remainingPayment.Sub(paymentToRegularDebt).RoundBank(config.Round)
}

func (period *LoanPeriod) applyPaymentToExtraPrincipal(periodMovement *LoanPeriodMovement, payment decimal.Decimal) decimal.Decimal {
	remainingPayment := payment
	if remainingPayment.LessThanOrEqual(decimal.Zero) {
		return remainingPayment
	}
	paymentToExtraPrincipal := remainingPayment
	periodMovement.PaidExtraToPrincipal = paymentToExtraPrincipal
	period.TotalPaidExtraToPrincipal = period.TotalPaidExtraToPrincipal.Add(paymentToExtraPrincipal).RoundBank(config.Round)
	return remainingPayment.Sub(paymentToExtraPrincipal).RoundBank(config.Round)
}

func calculateDaysLate(lastLiquidationDate, liquidationDate time.Time) int {
	daysLate := 0
	if liquidationDate.After(lastLiquidationDate) {
		daysLate = utils.DaysBetween(lastLiquidationDate, liquidationDate)
	}
	return daysLate
}
