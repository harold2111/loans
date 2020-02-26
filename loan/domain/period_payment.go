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
	ID          int `gorm:"primary_key"`
	PeriodID    int
	PaymentID   int
	Amount      decimal.Decimal `gorm:"not null; type:numeric"`
	PaymentType string
}

func newPeriodPayment(periodID int, paymentID int, amount decimal.Decimal, paymentType string) PeriodPayment {
	return PeriodPayment{0, periodID, paymentID, amount, paymentType}
}
