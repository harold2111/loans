package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	LoanID      uint            `json:"loanID"`
	Payment     decimal.Decimal `json:"payment"`
	PaymentDate time.Time       `json:"paymentDate"`
}
