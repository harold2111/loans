package domain

import "github.com/shopspring/decimal"

//DefaultPayment represents payment to a default
type DefaultPayment struct {
	DefaultID int
	PaymentID int
	Amount    decimal.Decimal
}

func newDefaultPayment(defaultID int, paymentID int, amount decimal.Decimal) DefaultPayment {
	return DefaultPayment{defaultID, paymentID, amount}
}
