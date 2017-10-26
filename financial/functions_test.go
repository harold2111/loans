package financial

import (
	"testing"

	"github.com/shopspring/decimal"
)

const (
	round int32 = 6
)

func TestCalculatePaymentExpectedSuccess(t *testing.T) {
	principal := decimal.NewFromFloat(5000000)
	interestRatePeriod := decimal.NewFromFloat(0.03)
	periodNumbers := 24
	paymentExpected := decimal.NewFromFloat(295237.079745)
	paymentActual := CalculatePayment(principal, interestRatePeriod, periodNumbers).RoundBank(round)
	if !paymentActual.Equal(paymentExpected) {
		t.Fatalf("Expected %s but got %s", paymentExpected, paymentActual)
	}
}

func TestBalanceExpectedInSpecificPeriodExpectedSuccess(t *testing.T) {
	principal := decimal.NewFromFloat(5000000)
	interestRatePeriod := decimal.NewFromFloat(0.03)
	periodNumbers := 24
	specificPeriod := 17
	balance := BalanceExpectedInSpecificPeriod(principal, interestRatePeriod, periodNumbers, specificPeriod)

	balanceExpected := decimal.NewFromFloat(1839410.545684)
	balanceActual := balance.FinalBalance.RoundBank(round)
	if !balanceActual.Equal(balanceExpected) {
		t.Fatalf("Expected %s but got %s", balanceExpected, balanceActual)
	}

	toPrincipalExpect := decimal.NewFromFloat(233062.877062)
	toPrincipalActual := balance.ToPrincipal.RoundBank(round)
	if !toPrincipalActual.Equal(toPrincipalExpect) {
		t.Fatalf("Expected %s but got %s", toPrincipalExpect, toPrincipalActual)
	}
}

func TestEffectiveMonthlyToAnnualExpectedSuccess(t *testing.T) {
	annualExpected := decimal.NewFromFloat(0.268242)
	monthly := decimal.NewFromFloat(0.02)
	annualActual := EffectiveMonthlyToAnnual(monthly).RoundBank(round)

	if !annualActual.Equal(annualExpected) {
		t.Fatalf("Expected %s but got %s", annualExpected, annualActual)
	}
}

func TestCalculateInterestPastOfDueExpectedSuccess(t *testing.T) {
	effectiveAnnualInterestRateForLate := decimal.NewFromFloat(0.3022)
	due := decimal.NewFromFloat(5000000)
	daysLate := 16
	interestPastDueExpected := decimal.NewFromFloat(66054.644809)
	interestPastDueActual := CalculateInterestPastOfDue(effectiveAnnualInterestRateForLate, due, daysLate).RoundBank(round)

	if !interestPastDueActual.Equal(interestPastDueExpected) {
		t.Fatalf("Expected %s but got %s", interestPastDueExpected, interestPastDueActual)
	}
}
