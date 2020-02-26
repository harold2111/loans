package application

import (
	"time"

	"github.com/shopspring/decimal"
)

type LoanResponse struct {
	ID                 int             `json:"id"`
	Principal          decimal.Decimal `json:"principal"`
	InterestRatePeriod decimal.Decimal `json:"interestRatePeriod"`
	PeriodNumbers      int             `json:"periodNumbers"`
	StartDate          time.Time       `json:"startDate"`
	ClientID           int             `json:"clientID"`
	PaymentAgreed      decimal.Decimal `json:"paymentAgreed"`
	CloseDateAgreed    time.Time       `json:"closeDateAgreed"`
	State              string          `json:"status"`
}
