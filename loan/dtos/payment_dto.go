package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

type PaymentDTO struct {
	ID            uint            `gorm:"primary_key" json:"id"`
	LoanID        uint            `gorm:"not null" json:"loanID" validate:"required"`
	PaymentAmount decimal.Decimal `gorm:"not null; type:numeric" json:"paymentAmount"`
	PaymentDate   time.Time       `json:"paymentDate"`
}
