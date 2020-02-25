package domain

import (
	"time"

	"github.com/harold2111/loans/shared/config"
	"github.com/shopspring/decimal"
)

//DefaultPeriod represents period defaults
type DefaultPeriod struct {
	ID              int
	PeriodID        int
	LiquidationDate time.Time
	DaysInDefault   int
	DebtForDefault  decimal.Decimal `gorm:"type:numeric"`
	PaidToDefault   decimal.Decimal `gorm:"type:numeric"`
	Payments        []DefaultPayment
}

func newDefaultPeriod(periodID int, liquidationDate time.Time, daysInDefault int, debtForDefault decimal.Decimal) DefaultPeriod {
	return DefaultPeriod{0, periodID, liquidationDate, daysInDefault, debtForDefault, decimal.Zero, []DefaultPayment{}}
}

func (d *DefaultPeriod) applyPayment(paymentID int, paymentAmount decimal.Decimal) decimal.Decimal {
	remainingPayment := paymentAmount
	if remainingPayment.LessThanOrEqual(decimal.Zero) || d.totalDebt().LessThanOrEqual(decimal.Zero) {
		return paymentAmount
	}
	var paymentToDefault decimal.Decimal
	if remainingPayment.LessThanOrEqual(d.totalDebt()) {
		paymentToDefault = remainingPayment
	} else {
		paymentToDefault = d.totalDebt()
	}
	d.PaidToDefault = d.PaidToDefault.Add(paymentToDefault).RoundBank(config.Round)
	defaultPayment := newDefaultPayment(d.ID, paymentID, paymentToDefault)
	d.Payments = append(d.Payments, defaultPayment)
	remainingPayment = remainingPayment.Sub(paymentToDefault)
	return remainingPayment.RoundBank(config.Round)
}

func (d DefaultPeriod) totalDebt() decimal.Decimal {
	return d.DebtForDefault.Sub(d.PaidToDefault)
}
