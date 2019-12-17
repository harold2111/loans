package financial

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	InitialPrincipal   decimal.Decimal `json:"initialPrincipal"`
	Payment            decimal.Decimal `json:"payment"`
	InterestRatePeriod decimal.Decimal `json:"interestRatePeriod"`
	ToInterest         decimal.Decimal `json:"toInterest"`
	ToPrincipal        decimal.Decimal `json:"toPrincipal"`
	FinalPrincipal     decimal.Decimal `json:"finalPrincipal"`
}

func (balance *Balance) calculateAmountBalance() {
	//round := config.Round
	toInterest := balance.InitialPrincipal.Mul(balance.InterestRatePeriod)
	toPrincipal := balance.Payment.Sub(toInterest)
	finalPrincipal := balance.InitialPrincipal.Sub(toPrincipal)

	balance.ToInterest = toInterest
	balance.ToPrincipal = toPrincipal
	balance.FinalPrincipal = finalPrincipal
}
