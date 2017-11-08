package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	LoanID        uint            `json:"loanID"`
	PaymentAmount decimal.Decimal `json:"paymentAmount"`
	PaymentDate   time.Time       `json:"paymentDate"`
}

type PaymentResponse struct {
	ID uint
	Payment
}
