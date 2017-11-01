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
	nextBalance := Balance{}
	nextBalance.InitialPrincipal = beforeBalance.FinalPrincipal
	nextBalance.Payment = beforeBalance.Payment
	nextBalance.InterestRatePeriod = beforeBalance.InterestRatePeriod
	nextBalance.calculateAmountBalance()
	return nextBalance
}

func (balance *Balance) calculateAmountBalance() {
	balance.ToInterest = balance.InitialPrincipal.Mul(balance.InterestRatePeriod)
	balance.ToPrincipal = balance.Payment.Sub(balance.ToInterest)
	balance.FinalPrincipal = balance.InitialPrincipal.Sub(balance.ToPrincipal)
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
