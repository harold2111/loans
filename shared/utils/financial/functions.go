package financial

import (
	"loans/shared/config"
	"strconv"

	"github.com/shopspring/decimal"
)

func CalculatePayment(principal decimal.Decimal, interestRatePeriod decimal.Decimal, periodNumbers int) decimal.Decimal {
	one := decimal.NewFromFloat(1)
	n, _ := decimal.NewFromString(strconv.Itoa(periodNumbers))
	nNeg := n.Neg()
	rate := interestRatePeriod

	rateMulPrincipal := rate.Mul(principal)
	ratePlusOne := rate.Add(one)
	ratePlusOnePowNNeg := ratePlusOne.Pow(nNeg)
	oneMinusRatePlusOnePowNNeg := one.Sub(ratePlusOnePowNNeg)

	payment := rateMulPrincipal.Div(oneMinusRatePlusOnePowNNeg)

	return payment
}

func Amortizations(principal decimal.Decimal, interestRatePeriod decimal.Decimal, periodNumbers int) []Balance {
	round := config.Round
	balances := make([]Balance, periodNumbers)
	balances[0].InitialPrincipal = principal.Truncate(round)
	balances[0].Payment = CalculatePayment(principal, interestRatePeriod, periodNumbers).Truncate(round)
	balances[0].InterestRatePeriod = interestRatePeriod.Truncate(round)
	balances[0].calculateAmountBalance()
	for period := 1; period < periodNumbers; period++ {
		balances[period] = NextBalanceFromBefore(balances[period-1])
	}
	return balances
}

func BalanceExpectedInSpecificPeriod(principal decimal.Decimal, interestRatePeriod decimal.Decimal, periodNumbers int, specificPeriod int) Balance {
	payment := CalculatePayment(principal, interestRatePeriod, periodNumbers)
	initialBalance := Balance{}
	initialBalance.InitialPrincipal = principal
	initialBalance.Payment = payment
	initialBalance.InterestRatePeriod = interestRatePeriod
	initialBalance.calculateAmountBalance()
	finalBalance := initialBalance
	for period := 1; period < specificPeriod; period++ {
		finalBalance = NextBalanceFromBefore(initialBalance)
		initialBalance = finalBalance
	}
	return finalBalance
}

func NextBalanceFromBefore(beforeBalance Balance) Balance {
	round := config.Round
	nextBalance := Balance{}
	nextBalance.InitialPrincipal = beforeBalance.FinalPrincipal.Truncate(round)
	nextBalance.Payment = beforeBalance.Payment.Truncate(round)
	nextBalance.InterestRatePeriod = beforeBalance.InterestRatePeriod.Truncate(round)
	nextBalance.calculateAmountBalance()
	return nextBalance
}

func (balance *Balance) calculateAmountBalance() {
	round := config.Round
	balance.ToInterest = balance.InitialPrincipal.Mul(balance.InterestRatePeriod).Truncate(round)
	balance.ToPrincipal = balance.Payment.Sub(balance.ToInterest).Truncate(round)
	balance.FinalPrincipal = balance.InitialPrincipal.Sub(balance.ToPrincipal).Truncate(round)
}

func CalculateInterestPastOfDueDIAN(effectiveAnnualInterestRate, paymentLate decimal.Decimal, daysLate int) decimal.Decimal {
	dailyInterest := effectiveAnnualInterestRate.Div(decimal.NewFromFloat(366))
	return paymentLate.Mul(dailyInterest).Mul(decimal.NewFromFloat(float64(daysLate)))
}

func FeeLateWithPeriodInterest(periodInteres, paymentLate decimal.Decimal, daysLate int) decimal.Decimal {
	dailyInterest := periodInteres.Div(decimal.NewFromFloat(30))
	return paymentLate.Mul(dailyInterest).Mul(decimal.NewFromFloat(float64(daysLate)))
}

func EffectiveMonthlyToAnnual(monthlyRate decimal.Decimal) decimal.Decimal {
	one := decimal.NewFromFloat(1)
	onePlusMonthlyRate := one.Add(monthlyRate)
	onePlusMonthlyRatePoWTwelve := onePlusMonthlyRate.Pow(decimal.NewFromFloat(12))
	onePlusMonthlyRatePoWTwelveMinusOne := onePlusMonthlyRatePoWTwelve.Sub(one)
	return onePlusMonthlyRatePoWTwelveMinusOne
}
