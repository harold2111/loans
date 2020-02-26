package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

const (

	//ExtraToPrincipal if the payment has extra amount should be applied to the principal.
	ExtraToPrincipal = "ExtraToPrincipal"
	//ExtraToNextPeriods if the payment has extra amount should be applied to the next periods.
	ExtraToNextPeriods = "ExtraToNextPeriods"
)

//Payment represents a loan payment
type Payment struct {
	ID              int             `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time       `json:"-"`
	UpdatedAt       time.Time       `json:"-"`
	DeletedAt       *time.Time      `sql:"index" json:"-"`
	LoanID          int             `gorm:"not null" json:"loanID" validate:"required"`
	PaymentAmount   decimal.Decimal `gorm:"not null; type:numeric" json:"paymentAmount"`
	RemainingAmount decimal.Decimal `gorm:"not null; type:numeric" json:"remainingPayment"`
	PaymentDate     time.Time       `json:"paymentDate"`
	PaymentType     string          `json:"paymentType"`
}

//NewPayment create a new payment
func NewPayment(loanID int, paymentAmount decimal.Decimal, paymentDate time.Time, paymentType string) Payment {
	return Payment{
		LoanID:          loanID,
		PaymentAmount:   paymentAmount,
		RemainingAmount: paymentAmount,
		PaymentDate:     paymentDate,
		PaymentType:     paymentType,
	}
}

func (p Payment) isExtraToPrincipal() bool {
	return p.PaymentType == ExtraToPrincipal
}

func (p Payment) isExtraToNextPeriods() bool {
	return p.PaymentType == ExtraToNextPeriods
}
