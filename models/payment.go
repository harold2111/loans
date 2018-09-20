package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	ID            uint            `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	DeletedAt     *time.Time      `sql:"index" json:"deletedAt"`
	LoanID        uint            `gorm:"not null" json:"loanID" validate:"required"`
	PaymentAmount decimal.Decimal `gorm:"not null; type:numeric" json:"paymentAmount"`
	PaymentDate   time.Time       `json:"paymentDate"`
}
