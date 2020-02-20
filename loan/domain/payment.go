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
	ID            uint            `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time       `json:"-"`
	UpdatedAt     time.Time       `json:"-"`
	DeletedAt     *time.Time      `sql:"index" json:"-"`
	LoanID        uint            `gorm:"not null" json:"loanID" validate:"required"`
	PaymentAmount decimal.Decimal `gorm:"not null; type:numeric" json:"paymentAmount"`
	PaymentDate   time.Time       `json:"paymentDate"`
	PaymentType   string          `json:"paymentType"`
}
