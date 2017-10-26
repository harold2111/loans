package financial

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	Payment        decimal.Decimal
	InitialBalance decimal.Decimal
	ToInterest     decimal.Decimal
	ToPrincipal    decimal.Decimal
	FinalBalance   decimal.Decimal
}
