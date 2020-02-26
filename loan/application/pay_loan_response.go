package application

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayLoanResponse struct {
	ID              int             `gorm:"primary_key" json:"id"`
	LoanID          int             `gorm:"not null" json:"loanID" validate:"required"`
	PaymentAmount   decimal.Decimal `gorm:"not null; type:numeric" json:"paymentAmount"`
	RemainingAmount decimal.Decimal `gorm:"not null; type:numeric" json:"remainingAmount"`
	PaymentDate     time.Time       `json:"paymentDate"`
}
