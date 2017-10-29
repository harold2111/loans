package financial

import (
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

func BalanceExpectedInSpecificPeriod(principal decimal.Decimal, interestRatePeriod decimal.Decimal, periodNumbers int, specificPeriod int) Balance {
	payment := CalculatePayment(principal, interestRatePeriod, periodNumbers)
	initialBalance := principal
	toInterest := decimal.Zero
	toPrincipal := decimal.Zero
	finalBalance := decimal.Zero
	for period := 0; period < specificPeriod; period++ {
		if period > 0 {
			initialBalance = finalBalance
		}
		toInterest = initialBalance.Mul(interestRatePeriod)
		toPrincipal = payment.Sub(toInterest)
		finalBalance = initialBalance.Sub(toPrincipal)

	}
	return Balance{
		Payment:        payment,
		InitialBalance: initialBalance,
		ToInterest:     toInterest,
		ToPrincipal:    toPrincipal,
		FinalBalance:   finalBalance}
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
