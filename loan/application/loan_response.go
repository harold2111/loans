package application

import (
	"time"

	"github.com/shopspring/decimal"
)

type LoanResponse struct {
	ID                 uint            `json:"id"`
	Principal          decimal.Decimal `json:"principal"`
	InterestRatePeriod decimal.Decimal `json:"interestRatePeriod"`
	PeriodNumbers      uint            `json:"periodNumbers"`
	StartDate          time.Time       `json:"startDate"`
	ClientID           uint            `json:"clientID"`
	PaymentAgreed      decimal.Decimal `json:"paymentAgreed"`
	CloseDateAgreed    time.Time       `json:"closeDateAgreed"`
	State              string          `json:"status"`
}
