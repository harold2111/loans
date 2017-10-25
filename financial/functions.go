package financial

import (
	"fmt"
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

func BalancingExpectedInSpecificPeriodNumber(principal decimal.Decimal, interestRatePeriod decimal.Decimal, periodNumbers int, specificPeriod int) decimal.Decimal {
	payment := CalculatePayment(principal, interestRatePeriod, periodNumbers)
	initialBalance := principal
	var finalBalance decimal.Decimal
	for period := 0; period < specificPeriod; period++ {
		toInterest := initialBalance.Mul(interestRatePeriod)
		toCapital := payment.Sub(toInterest)
		finalBalance = initialBalance.Sub(toCapital)
		initialBalance = finalBalance
	}
	return finalBalance
}

func CalculateInterestPastOfDue(effectiveAnnualInterestRate, paymentLate decimal.Decimal, daysLate int) decimal.Decimal {
	dailyInterest := effectiveAnnualInterestRate.Div(decimal.NewFromFloat(366))
	fmt.Println("dailyInterest", dailyInterest)
	return paymentLate.Mul(dailyInterest).Mul(decimal.NewFromFloat(float64(daysLate)))
}

func EffectiveMonthlyToAnnual(monthlyRate decimal.Decimal) decimal.Decimal {
	one := decimal.NewFromFloat(1)
	onePlusMonthlyRate := one.Add(monthlyRate)
	onePlusMonthlyRatePoWTwelve := onePlusMonthlyRate.Pow(decimal.NewFromFloat(12))
	onePlusMonthlyRatePoWTwelveMinusOne := onePlusMonthlyRatePoWTwelve.Sub(one)
	return onePlusMonthlyRatePoWTwelveMinusOne
}
