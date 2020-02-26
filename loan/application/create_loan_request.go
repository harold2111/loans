package application

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateLoanRequest struct {
	Principal          decimal.Decimal `json:"principal"`
	InterestRatePeriod decimal.Decimal `json:"interestRatePeriod"`
	PeriodNumbers      int             `json:"periodNumbers"`
	StartDate          time.Time       `json:"startDate"`
	ClientID           int             `json:"clientID"`
}
