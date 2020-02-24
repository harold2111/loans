package domain

import (
	"sort"
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
	ID                 int `gorm:"primary_key"`
	LoanID             int
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time `sql:"index"`
	PeriodNumber       int
	State              string
	StartDate          time.Time
	EndDate            time.Time
	MaxPaymentDate     time.Time
	InitialPrincipal   decimal.Decimal `gorm:"type:numeric"`
	Payment            decimal.Decimal `gorm:"type:numeric"`
	InterestRate       decimal.Decimal `gorm:"type:numeric"`
	PrincipalOfPayment decimal.Decimal `gorm:"type:numeric"`
	InterestOfPayment  decimal.Decimal `gorm:"type:numeric"`
	FinalPrincipal     decimal.Decimal `gorm:"type:numeric"`
	//Modifiable fields
	PeriodPayments            []PeriodPayment
	PeriodDefaults            []PeriodDefault
	TotalPaidToRegularDebt    decimal.Decimal `gorm:"type:numeric"`
	TotalPaidExtraToPrincipal decimal.Decimal `gorm:"type:numeric"`
}

func (period *LoanPeriod) liquidateByDate(liquidationDate time.Time) {
	if period.isPeriodOnDebt(liquidationDate) {
		period.State = LoanPeriodStateDue
	}
	if period.hasToCalculateDebtForDefaults(liquidationDate) {
		period.calculateDebtForDefaults(liquidationDate)
	}
}

func (period *LoanPeriod) applyRegularPayment(payment Payment) Payment {
	remainingPayment := payment.PaymentAmount
	remainingPayment = period.applyPaymentToDefaults(payment.ID, remainingPayment)
	remainingPayment = period.applyPaymentToRegularDebt(payment.ID, remainingPayment)
	if period.TotalDebt().LessThanOrEqual(decimal.Zero) {
		period.State = LoanPeriodStatePaid
	}
	payment.PaymentAmount = remainingPayment
	return payment
}

func (period *LoanPeriod) applyPaymentToDefaults(paymentID int, paymentAmount decimal.Decimal) decimal.Decimal {
	remainingPayment := paymentAmount
	totalDefaultDebt := period.TotalDefaultDebt()
	if totalDefaultDebt.LessThanOrEqual(decimal.Zero) || remainingPayment.LessThanOrEqual(decimal.Zero) {
		return remainingPayment
	}
	for periodDefaultIndex := 0; periodDefaultIndex < len(period.PeriodDefaults); periodDefaultIndex++ {
		remainingPayment = period.PeriodDefaults[periodDefaultIndex].applyPayment(paymentID, paymentAmount)

	}
	return remainingPayment.RoundBank(config.Round)
}

func (period *LoanPeriod) applyPaymentToRegularDebt(paymentID int, payment decimal.Decimal) decimal.Decimal {
	remainingPayment := payment
	totalDebtOfPayment := period.TotalRegularDebt()
	if totalDebtOfPayment.LessThanOrEqual(decimal.Zero) || remainingPayment.LessThanOrEqual(decimal.Zero) {
		return remainingPayment
	}
	var paymentToRegularDebt decimal.Decimal
	if remainingPayment.LessThanOrEqual(totalDebtOfPayment) {
		paymentToRegularDebt = remainingPayment
	} else {
		paymentToRegularDebt = totalDebtOfPayment
	}
	periodPayment := newPeriodPayment(period.ID, paymentID, paymentToRegularDebt, PaymentTypeRegular)
	period.PeriodPayments = append(period.PeriodPayments, periodPayment)
	period.TotalPaidToRegularDebt = period.TotalPaidToRegularDebt.Add(paymentToRegularDebt).RoundBank(config.Round)
	return remainingPayment.Sub(paymentToRegularDebt).RoundBank(config.Round)
}

func (period *LoanPeriod) applyPaymentToPrincipal(paymentID int, payment Payment) Payment {
	remainingPayment := payment
	if remainingPayment.PaymentAmount.LessThanOrEqual(decimal.Zero) {
		return remainingPayment
	}
	var paymentToExtraPrincipal decimal.Decimal
	if remainingPayment.PaymentAmount.GreaterThanOrEqual(period.FinalPrincipal) {
		paymentToExtraPrincipal = period.FinalPrincipal
	} else {
		paymentToExtraPrincipal = remainingPayment.PaymentAmount
	}
	periodPayment := newPeriodPayment(period.ID, paymentID, paymentToExtraPrincipal, PaymentTypePrincipal)
	period.PeriodPayments = append(period.PeriodPayments, periodPayment)
	period.FinalPrincipal = period.FinalPrincipal.Sub(paymentToExtraPrincipal).RoundBank(config.Round)
	period.TotalPaidExtraToPrincipal = period.TotalPaidExtraToPrincipal.Add(paymentToExtraPrincipal).RoundBank(config.Round)
	remainingPayment.PaymentAmount = remainingPayment.PaymentAmount.Sub(paymentToExtraPrincipal).RoundBank(config.Round)
	return remainingPayment
}

func (period *LoanPeriod) calculateDebtForDefaults(liquidationDate time.Time) {
	lastLiquidationDate := period.lastLiquidationDate()
	daysInDefaultSinceLastLiquidation := calculateDaysLate(lastLiquidationDate, liquidationDate)
	if daysInDefaultSinceLastLiquidation > 0 {
		daysInDefault := daysInDefaultSinceLastLiquidation
		debtForDefault := financial.FeeLateWithPeriodInterest(period.InterestRate, period.TotalRegularDebt(), daysInDefaultSinceLastLiquidation).RoundBank(config.Round)
		periodDefault := newPeriodDefault(period.ID, liquidationDate, daysInDefault, debtForDefault)
		period.PeriodDefaults = append(period.PeriodDefaults, periodDefault)
	}
}

//TotalRegularDebt returns the total regular payment debt of period
func (period LoanPeriod) TotalRegularDebt() decimal.Decimal {
	if period.State == LoanPeriodStateAnnulled {
		return decimal.Zero
	}
	return period.Payment.Sub(period.TotalPaidToRegularDebt).RoundBank(config.Round)
}

//TotalDefaultDebt returns the total default debt of period
func (period LoanPeriod) TotalDefaultDebt() decimal.Decimal {
	if period.State == LoanPeriodStateAnnulled {
		return decimal.Zero
	}
	var totalDefaultDebt decimal.Decimal
	for _, periodDefault := range period.PeriodDefaults {
		totalDefaultDebt = totalDefaultDebt.Add(periodDefault.totalDebt())
	}
	return totalDefaultDebt
}

//TotalDaysInDefault returns the total days in default of all periods
func (period LoanPeriod) TotalDaysInDefault() int {
	var totalDaysInDefault int
	for _, periodDefault := range period.PeriodDefaults {
		totalDaysInDefault = totalDaysInDefault + periodDefault.DaysInDefault
	}
	return totalDaysInDefault
}

//TotalDebt total debt of period
func (period LoanPeriod) TotalDebt() decimal.Decimal {
	totalRegularDebt := period.TotalRegularDebt()
	totalDefaultDebt := period.TotalDefaultDebt()
	return totalDefaultDebt.Add(totalRegularDebt)
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

func (period *LoanPeriod) hasToCalculateDebtForDefaults(liquidationDate time.Time) bool {
	return period.State == LoanPeriodStateDue && liquidationDate.After(period.MaxPaymentDate)
}

func (period *LoanPeriod) lastLiquidationDate() time.Time {
	periodDefaults := period.PeriodDefaults
	var lastLiquidationDate time.Time
	if len(periodDefaults) > 0 {
		sort.Slice(periodDefaults, func(i, j int) bool {
			return periodDefaults[i].LiquidationDate.Before(periodDefaults[j].LiquidationDate)
		})
		lastLiquidationDate = periodDefaults[len(periodDefaults)-1].LiquidationDate
	} else {
		lastLiquidationDate = period.EndDate
	}
	return lastLiquidationDate
}
