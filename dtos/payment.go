package dtos

import (
	"github.com/shopspring/decimal"
)

type Payment struct {
	LoanID  uint            `json:"loanID"`
	Payment decimal.Decimal `json:"payment"`
}
