package domain

import "github.com/shopspring/decimal"

const (
	//PaymentTypeRegular represents a regular payment
	PaymentTypeRegular = "REGULAR"
	//PaymentTypePrincipal represents a extra payment to the principal
	PaymentTypePrincipal = "PRINCIPAL"
)

//PeriodPayment represents payment to a default
type PeriodPayment struct {
	PeriodID    int
	PaymentID   int
	Amount      decimal.Decimal
	PaymentType string
}

func newPeriodPayment(periodID int, paymentID int, amount decimal.Decimal, paymentType string) PeriodPayment {
	return PeriodPayment{periodID, paymentID, amount, paymentType}
}
