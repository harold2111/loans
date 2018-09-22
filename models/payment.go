package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	ID            uint            `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time       `json:"-"`
	UpdatedAt     time.Time       `json:"-"`
	DeletedAt     *time.Time      `sql:"index" json:"-"`
	LoanID        uint            `gorm:"not null" json:"loanID" validate:"required"`
	PaymentAmount decimal.Decimal `gorm:"not null; type:numeric" json:"paymentAmount"`
	PaymentDate   time.Time       `json:"paymentDate"`
}
