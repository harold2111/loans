package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateLoanRequest struct {
	Principal          decimal.Decimal `json:"principal"`
	InterestRatePeriod decimal.Decimal `json:"interestRatePeriod"`
	PeriodNumbers      uint            `json:"periodNumbers"`
	StartDate          time.Time       `json:"startDate"`
	ClientID           uint            `json:"clientID"`
}
