package domain

import "github.com/shopspring/decimal"

//DefaultPayment represents payment to a default
type DefaultPayment struct {
	ID              int `gorm:"primary_key"`
	DefaultPeriodID int
	PaymentID       int
	Amount          decimal.Decimal `gorm:"not null; type:numeric"`
}

func newDefaultPayment(defaultPeriodID int, paymentID int, amount decimal.Decimal) DefaultPayment {
	return DefaultPayment{0, defaultPeriodID, paymentID, amount}
}
