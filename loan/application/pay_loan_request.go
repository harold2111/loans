package application

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayLoanRequest struct {
	LoanID        int             `json:"loanID"`
	PaymentAmount decimal.Decimal `json:"paymentAmount"`
	PaymentDate   time.Time       `json:"paymentDate"`
	PaymentType   string          `json:"paymentType"`
}
