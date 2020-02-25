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
	//PeriodStateOpen period open.
	PeriodStateOpen = "OPEN"
	//PeriodStateDue period state due.
	PeriodStateDue = "DUE"
	//PeriodStatePaid period state paid.
	PeriodStatePaid = "PAID"
	//PeriodStateAnnulled period state annulled.
	PeriodStateAnnulled = "ANNULLED"
)

//Period represents a period of a loan.
type Period struct {
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
	Payments                  []PeriodPayment
	DefaultPeriods            []DefaultPeriod
	TotalPaidToRegularDebt    decimal.Decimal `gorm:"type:numeric"`
	TotalPaidExtraToPrincipal decimal.Decimal `gorm:"type:numeric"`
}

func (period *Period) liquidateByDate(liquidationDate time.Time) {
	if period.isPeriodOnDebt(liquidationDate) {
		period.State = PeriodStateDue
	}
	if period.hasToCalculateDebtForDefaults(liquidationDate) {
		period.calculateDebtForDefaults(liquidationDate)
	}
}

func (period *Period) applyRegularPayment(payment Payment) Payment {
	remainingPayment := payment.PaymentAmount
	remainingPayment = period.applyPaymentToDefaults(payment.ID, remainingPayment)
	remainingPayment = period.applyPaymentToRegularDebt(payment.ID, remainingPayment)
	if period.TotalDebt().LessThanOrEqual(decimal.Zero) {
		period.State = PeriodStatePaid
	}
	payment.PaymentAmount = remainingPayment
	return payment
}

func (period *Period) applyPaymentToDefaults(paymentID int, paymentAmount decimal.Decimal) decimal.Decimal {
	remainingPayment := paymentAmount
	totalDefaultDebt := period.TotalDefaultDebt()
	if totalDefaultDebt.LessThanOrEqual(decimal.Zero) || remainingPayment.LessThanOrEqual(decimal.Zero) {
		return remainingPayment
	}
	for periodDefaultIndex := 0; periodDefaultIndex < len(period.DefaultPeriods); periodDefaultIndex++ {
		remainingPayment = period.DefaultPeriods[periodDefaultIndex].applyPayment(paymentID, paymentAmount)
	}
	return remainingPayment.RoundBank(config.Round)
}

func (period *Period) applyPaymentToRegularDebt(paymentID int, payment decimal.Decimal) decimal.Decimal {
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
	period.Payments = append(period.Payments, periodPayment)
	period.TotalPaidToRegularDebt = period.TotalPaidToRegularDebt.Add(paymentToRegularDebt).RoundBank(config.Round)
	return remainingPayment.Sub(paymentToRegularDebt).RoundBank(config.Round)
}

func (period *Period) applyPaymentToPrincipal(paymentID int, payment Payment) Payment {
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
	period.Payments = append(period.Payments, periodPayment)
	period.FinalPrincipal = period.FinalPrincipal.Sub(paymentToExtraPrincipal).RoundBank(config.Round)
	period.TotalPaidExtraToPrincipal = period.TotalPaidExtraToPrincipal.Add(paymentToExtraPrincipal).RoundBank(config.Round)
	remainingPayment.PaymentAmount = remainingPayment.PaymentAmount.Sub(paymentToExtraPrincipal).RoundBank(config.Round)
	return remainingPayment
}

func (period *Period) calculateDebtForDefaults(liquidationDate time.Time) {
	daysInDefaultSinceLastLiquidation := period.calculateDaysLate(liquidationDate)
	if daysInDefaultSinceLastLiquidation > 0 {
		daysInDefault := daysInDefaultSinceLastLiquidation
		debtForDefault := financial.FeeLateWithPeriodInterest(period.InterestRate, period.TotalRegularDebt(), daysInDefaultSinceLastLiquidation).RoundBank(config.Round)
		periodDefault := newDefaultPeriod(period.ID, liquidationDate, daysInDefault, debtForDefault)
		period.DefaultPeriods = append(period.DefaultPeriods, periodDefault)
	}
}

func (period *Period) annullate() {
	period.State = PeriodStateAnnulled
}

//TotalRegularDebt returns the total regular payment debt of period
func (period Period) TotalRegularDebt() decimal.Decimal {
	if period.State == PeriodStateAnnulled {
		return decimal.Zero
	}
	return period.Payment.Sub(period.TotalPaidToRegularDebt).RoundBank(config.Round)
}

//TotalDefaultDebt returns the total default debt of period
func (period Period) TotalDefaultDebt() decimal.Decimal {
	if period.State == PeriodStateAnnulled {
		return decimal.Zero
	}
	var totalDefaultDebt decimal.Decimal
	for _, periodDefault := range period.DefaultPeriods {
		totalDefaultDebt = totalDefaultDebt.Add(periodDefault.totalDebt())
	}
	return totalDefaultDebt
}

//TotalDaysInDefault returns the total days in default of all periods
func (period Period) TotalDaysInDefault() int {
	var totalDaysInDefault int
	for _, periodDefault := range period.DefaultPeriods {
		totalDaysInDefault = totalDaysInDefault + periodDefault.DaysInDefault
	}
	return totalDaysInDefault
}

//TotalDebt total debt of period
func (period Period) TotalDebt() decimal.Decimal {
	totalRegularDebt := period.TotalRegularDebt()
	totalDefaultDebt := period.TotalDefaultDebt()
	return totalDefaultDebt.Add(totalRegularDebt)
}

func (period Period) calculateDaysLate(liquidationDate time.Time) int {
	lastLiquidationDate := period.lastLiquidationDate()
	daysLate := 0
	if liquidationDate.After(lastLiquidationDate) {
		daysLate = utils.DaysBetween(lastLiquidationDate, liquidationDate)
	}
	return daysLate
}

func (period *Period) isPeriodOnDebt(liquidationDate time.Time) bool {
	endate := period.EndDate.AddDate(0, 0, config.DaysBeforeEndDateToConsiderateDue)
	return period.State == PeriodStateOpen && (endate.Before(liquidationDate) || endate.Equal(liquidationDate))
}

func (period *Period) hasToCalculateDebtForDefaults(liquidationDate time.Time) bool {
	return period.State == PeriodStateDue && liquidationDate.After(period.MaxPaymentDate)
}

func (period *Period) lastLiquidationDate() time.Time {
	periodDefaults := period.DefaultPeriods
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
