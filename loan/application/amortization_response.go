package application

import (
	"time"

	"github.com/shopspring/decimal"
)

type AmortizationResponse struct {
	Period             int             `json:"period"`
	MaxPaymentDate     time.Time       `json:"paymentDate"`
	InitialPrincipal   decimal.Decimal `json:"initialPrincipal"`
	Payment            decimal.Decimal `json:"payment"`
	InterestRatePeriod decimal.Decimal `json:"interestRatePeriod"`
	ToInterest         decimal.Decimal `json:"toInterest"`
	ToPrincipal        decimal.Decimal `json:"toPrincipal"`
	FinalPrincipal     decimal.Decimal `json:"finalPrincipal"`
}
