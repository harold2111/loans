package financial

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	InitialPrincipal   decimal.Decimal
	Payment            decimal.Decimal
	InterestRatePeriod decimal.Decimal
	ToInterest         decimal.Decimal
	ToPrincipal        decimal.Decimal
	FinalPrincipal     decimal.Decimal
}
