package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreatePaymentRequest struct {
	LoanID        uint            `json:"loanID"`
	PaymentAmount decimal.Decimal `json:"paymentAmount"`
	PaymentDate   time.Time       `json:"paymentDate"`
}
