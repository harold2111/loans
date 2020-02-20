package domain

import (
	"time"

	"github.com/harold2111/loans/shared/config"
	"github.com/harold2111/loans/shared/utils"
	"github.com/harold2111/loans/shared/utils/financial"

	"github.com/shopspring/decimal"
)

const (
	//LoanPeriodStateOpen period open.
	LoanPeriodStateOpen = "OPEN"

	//LoanPeriodStateDue period state due.
	LoanPeriodStateDue = "DUE"

	//LoanPeriodStatePaid period state paid.
	LoanPeriodStatePaid = "PAID"

	//LoanPeriodStateAnnulled period state annulled.
	LoanPeriodStateAnnulled = "ANNULLED"
)

//LoanPeriod represents a period of a loan.
type LoanPeriod struct {
	ID                 uint `gorm:"primary_key"`
	LoanID             uint
	Movements          []LoanPeriodMovement
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
	//Modifiable fields
	LastPaymentDate                time.Time
	DaysInArrearsSinceLastPayment  int
	DebtForArrearsSinceLastPayment decimal.Decimal `gorm:"type:numeric"`
	TotalDaysInArrears             int
	TotalDebtForArrears            decimal.Decimal `gorm:"type:numeric"`
	TotalDebtOfPayment             decimal.Decimal `gorm:"type:numeric"`
	TotalDebt                      decimal.Decimal `gorm:"type:numeric"`
	TotalPaid                      decimal.Decimal `gorm:"type:numeric"`
	TotalPaidToDebtForArrears      decimal.Decimal `gorm:"type:numeric"`
	TotalPaidToRegularDebt         decimal.Decimal `gorm:"type:numeric"`
	TotalPaidExtraToPrincipal      decimal.Decimal `gorm:"type:numeric"`
}

func (period *LoanPeriod) liquidateByDate(liquidationDate time.Time) {
	if period.isPeriodOnDebt(liquidationDate) {
		period.State = LoanPeriodStateDue
	}
	liquidationDatePlusGraceDays := liquidationDate.AddDate(0, 0, config.DaysAfterEndDateToConsiderateInArrears)
	if period.hasToCalculateDebtForArrears(liquidationDatePlusGraceDays) {
		period.calculateDebtForArrear(liquidationDate)
	}
}

func (period *LoanPeriod) applyRegularPayment(payment Payment) Payment {
	var periodMovement LoanPeriodMovement
	periodMovement.fillInitialMovementFromPeriod(*period)
	totalPaymentReceived := payment.PaymentAmount
	periodMovement.PaymentID = payment.ID

	remainingPayment := payment.PaymentAmount
	remainingPayment = period.applyPaymentToDebtForArrears(&periodMovement, remainingPayment)
	remainingPayment = period.applyPaymentToDebtOfPayments(&periodMovement, remainingPayment)

	period.TotalDebt = period.TotalDebtForArrears.Add(period.TotalDebtOfPayment).RoundBank(config.Round)
	period.TotalPaid = period.TotalPaid.Add(totalPaymentReceived)
	if period.TotalDebt.LessThanOrEqual(decimal.Zero) {
		period.State = LoanPeriodStatePaid
	}
	periodMovement.fillFinalMovementFromPeriod(*period)
	period.Movements = append(period.Movements, periodMovement)
	payment.PaymentAmount = remainingPayment
	return payment
}

func (period *LoanPeriod) applyPaymentToPrincipal(payment Payment) Payment {
	remainingPayment := payment
	if remainingPayment.PaymentAmount.LessThanOrEqual(decimal.Zero) {
		return remainingPayment
	}
	var periodMovement LoanPeriodMovement
	periodMovement.fillInitialMovementFromPeriod(*period)
	periodMovement.PaymentID = payment.ID
	var paymentToExtraPrincipal decimal.Decimal
	if remainingPayment.PaymentAmount.GreaterThanOrEqual(period.FinalPrincipal) {
		paymentToExtraPrincipal = period.FinalPrincipal
	} else {
		paymentToExtraPrincipal = remainingPayment.PaymentAmount
	}
	periodMovement.PaidExtraToPrincipal = paymentToExtraPrincipal
	period.FinalPrincipal = period.FinalPrincipal.Sub(paymentToExtraPrincipal).RoundBank(config.Round)
	period.TotalPaidExtraToPrincipal = period.TotalPaidExtraToPrincipal.Add(paymentToExtraPrincipal).RoundBank(config.Round)
	periodMovement.fillFinalMovementFromPeriod(*period)
	period.Movements = append(period.Movements, periodMovement)
	remainingPayment.PaymentAmount = remainingPayment.PaymentAmount.Sub(paymentToExtraPrincipal).RoundBank(config.Round)
	return remainingPayment
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

func (period *LoanPeriod) calculateDebtForArrear(liquidationDate time.Time) {
	daysInArrearsSinceLastPayment := calculateDaysLate(period.LastPaymentDate, liquidationDate)
	if daysInArrearsSinceLastPayment > 0 {
		debtForArrearsSinceLastPayment := financial.FeeLateWithPeriodInterest(period.InterestRate, period.TotalDebtOfPayment, daysInArrearsSinceLastPayment)
		period.DaysInArrearsSinceLastPayment = daysInArrearsSinceLastPayment
		period.DebtForArrearsSinceLastPayment = debtForArrearsSinceLastPayment.RoundBank(config.Round)
		period.TotalDaysInArrears = period.TotalDaysInArrears + daysInArrearsSinceLastPayment
		period.TotalDebtForArrears = period.TotalDebtForArrears.Add(debtForArrearsSinceLastPayment).RoundBank(config.Round)
		period.TotalDebt = period.TotalDebtOfPayment.Add(period.TotalDebtForArrears).RoundBank(config.Round)
	}
}

func calculateDaysLate(lastPaymentDate, liquidationDate time.Time) int {
	daysLate := 0
	if liquidationDate.After(lastPaymentDate) {
		daysLate = utils.DaysBetween(lastPaymentDate, liquidationDate)
	}
	return daysLate
}

func (period *LoanPeriod) isPeriodOnDebt(liquidationDate time.Time) bool {
	endate := period.EndDate.AddDate(0, 0, config.DaysBeforeEndDateToConsiderateDue)
	return period.State == LoanPeriodStateOpen && (endate.Before(liquidationDate) || endate.Equal(liquidationDate))
}

func (period *LoanPeriod) hasToCalculateDebtForArrears(liquidationDatePlusGraceDays time.Time) bool {
	return period.State == LoanPeriodStateDue && liquidationDatePlusGraceDays.After(period.EndDate)
}
