package domain

import (
	"testing"
	"time"

	"github.com/harold2111/loans/shared/config"
	"github.com/harold2111/loans/shared/utils"
	"github.com/shopspring/decimal"
)

type createLoanArgs struct {
	principal          decimal.Decimal
	interestRatePeriod decimal.Decimal
	periodNumbers      uint
	startDate          time.Time
	clientID           uint
	liquidationDate    *time.Time
}

func TestNewLoanForCreate(t *testing.T) {
	type expected struct {
		paymentAgreed   decimal.Decimal
		closeDateAgreed time.Time
	}
	tests := []struct {
		name string
		args createLoanArgs
		want expected
	}{
		{
			"CreateLoanTest-1",
			createLoanArgs{toDecimal(450000.0), toDecimal(0.05), 36, toDate(2019, 12, 16), 1, nil},
			expected{toDecimal(27195.5057), toDate(2022, 12, 16)},
		},
		{
			"CreateLoanTest-2",
			createLoanArgs{toDecimal(1000.0), toDecimal(0.01), 12, toDate(2019, 12, 16), 1, nil},
			expected{toDecimal(88.8488), toDate(2020, 12, 16)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoanForCreate(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, config.DefaultGraceDays, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoanForCreate() error = %v", err)
				return
			}
			if got.ClientID != tt.args.clientID {
				t.Errorf("ClientID = %v, want %v", got.ClientID, tt.args.clientID)
			}
			if !got.Principal.Equal(tt.args.principal) {
				t.Errorf("Principal = %v, want %v", got.Principal, tt.args.principal)
			}
			if !got.InterestRatePeriod.Equal(tt.args.interestRatePeriod) {
				t.Errorf("InterestRatePeriod = %v, want %v", got.InterestRatePeriod, tt.args.interestRatePeriod)
			}
			if got.PeriodNumbers != tt.args.periodNumbers {
				t.Errorf("PeriodNumbers = %v, want %v", got.PeriodNumbers, tt.args.periodNumbers)
			}
			if !got.PaymentAgreed.Equal(tt.want.paymentAgreed) {
				t.Errorf("PaymentAgreed = %v, want %v", got.PaymentAgreed, tt.want.paymentAgreed)
			}
			if !got.StartDate.Equal(tt.args.startDate) {
				t.Errorf("StartDate = %v, want %v", got.StartDate, tt.args.startDate)
			}
			if got.State != LoanStateActive {
				t.Errorf("State = %v, want %v", got.State, LoanStateActive)
			}
			if !got.CloseDateAgreed.Equal(tt.want.closeDateAgreed) {
				t.Errorf("CloseDateAgreed = %v, want %v", got.CloseDateAgreed, tt.want.closeDateAgreed)
			}
			periodsLen := len(got.periods)
			if periodsLen != int(tt.args.periodNumbers) {
				t.Errorf("len(periods)  = %v, want %v", periodsLen, tt.args.periodNumbers)
			}
			lastPeriod := got.periods[periodsLen-1]
			if !lastPeriod.FinalPrincipal.Equal(decimal.Zero) {
				t.Errorf("lastPeriod.FinalPrincipal  = %v, want %v", lastPeriod.FinalPrincipal, decimal.Zero)
			}
		})
	}
}

func TestNewLoanForCreate_periods(t *testing.T) {
	type periodExpected struct {
		periodNumber       uint
		startDate          time.Time
		endDate            time.Time
		initialPrincipal   decimal.Decimal
		principalOfPayment decimal.Decimal
		interestOfPayment  decimal.Decimal
		finalPrincipal     decimal.Decimal
	}
	tests := []struct {
		name string
		args createLoanArgs
		want []periodExpected
	}{
		{
			"TestNewLoanForCreate_periods-1",
			createLoanArgs{toDecimal(3500000), toDecimal(0.02), 5, toDate(2019, 1, 31), 1, nil},
			[]periodExpected{
				{1, toDate(2019, 1, 31), toDate(2019, 2, 27), toDecimal(3500000), toDecimal(672554.3794), toDecimal(70000), toDecimal(2827445.6206)},
				{2, toDate(2019, 2, 28), toDate(2019, 3, 30), toDecimal(2827445.6206), toDecimal(686005.4670), toDecimal(56548.9124), toDecimal(2141440.1537)},
				{3, toDate(2019, 3, 31), toDate(2019, 4, 29), toDecimal(2141440.1537), toDecimal(699725.5763), toDecimal(42828.8031), toDecimal(1441714.5774)},
				{4, toDate(2019, 4, 30), toDate(2019, 5, 30), toDecimal(1441714.5774), toDecimal(713720.0878), toDecimal(28834.2915), toDecimal(727994.4896)},
				{5, toDate(2019, 5, 31), toDate(2019, 6, 29), toDecimal(727994.4896), toDecimal(727994.4896), toDecimal(14559.8898), toDecimal(0)},
			},
		},
		{
			"TestNewLoanForCreate_periods-2",
			createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1, nil},
			[]periodExpected{
				{1, toDate(2019, 1, 1), toDate(2019, 1, 31), toDecimal(100000), toDecimal(18648.1373), toDecimal(3500), toDecimal(81351.8627)},
				{2, toDate(2019, 2, 1), toDate(2019, 2, 28), toDecimal(81351.8627), toDecimal(19300.8221), toDecimal(2847.3152), toDecimal(62051.0406)},
				{3, toDate(2019, 3, 1), toDate(2019, 3, 31), toDecimal(62051.0406), toDecimal(19976.3509), toDecimal(2171.7864), toDecimal(42074.6897)},
				{4, toDate(2019, 4, 1), toDate(2019, 4, 30), toDecimal(42074.6897), toDecimal(20675.5232), toDecimal(1472.6141), toDecimal(21399.1665)},
				{5, toDate(2019, 5, 1), toDate(2019, 5, 31), toDecimal(21399.1665), toDecimal(21399.1665), toDecimal(748.9708), toDecimal(0)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoanForCreate(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, config.DefaultGraceDays, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoanForCreate() error = %v", err)
				return
			}
			for _, periodExpected := range tt.want {
				gotPeriod := got.periods[periodExpected.periodNumber-1]
				if gotPeriod.LoanID != got.ID {
					t.Errorf("LoanID = %v, want %v", gotPeriod.LoanID, got.ID)
				}
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if gotPeriod.State != LoanPeriodStateOpen {
					t.Errorf("State = %v, want %v", gotPeriod.PeriodNumber, LoanPeriodStateOpen)
				}
				if !gotPeriod.StartDate.Equal(periodExpected.startDate) {
					t.Errorf("StartDate = %v, want %v", gotPeriod.StartDate, periodExpected.startDate)
				}
				if !gotPeriod.EndDate.Equal(periodExpected.endDate) {
					t.Errorf("EndDate = %v, want %v", gotPeriod.EndDate, periodExpected.endDate)
				}
				if !gotPeriod.PaymentDate.Equal(periodExpected.endDate) {
					t.Errorf("PaymentDate = %v, want %v", gotPeriod.PaymentDate, periodExpected.endDate)
				}
				if !gotPeriod.InitialPrincipal.Equal(periodExpected.initialPrincipal) {
					t.Errorf("InitialPrincipal = %v, want %v", gotPeriod.InitialPrincipal, periodExpected.initialPrincipal)
				}
				if !gotPeriod.Payment.Equal(got.PaymentAgreed) {
					t.Errorf("Payment = %v, want %v", gotPeriod.Payment, got.PaymentAgreed)
				}
				if !gotPeriod.InterestRate.Equal(tt.args.interestRatePeriod) {
					t.Errorf("InterestRate = %v, want %v", gotPeriod.InterestRate, tt.args.interestRatePeriod)
				}
				if !gotPeriod.PrincipalOfPayment.Equal(periodExpected.principalOfPayment) {
					t.Errorf("PrincipalOfPayment = %v, want %v", gotPeriod.PrincipalOfPayment, periodExpected.principalOfPayment)
				}
				if !gotPeriod.InterestOfPayment.Equal(periodExpected.interestOfPayment) {
					t.Errorf("InterestOfPayment = %v, want %v", gotPeriod.InterestOfPayment, periodExpected.interestOfPayment)
				}
				if !gotPeriod.FinalPrincipal.Equal(periodExpected.finalPrincipal) {
					t.Errorf("FinalPrincipal = %v, want %v", gotPeriod.FinalPrincipal, periodExpected.finalPrincipal)
				}
				//Modifible fields
				if !gotPeriod.LastPaymentDate.Equal(periodExpected.endDate) {
					t.Errorf("LastPaymentDate = %v, want %v", gotPeriod.LastPaymentDate, periodExpected.endDate)
				}
				if gotPeriod.DaysInArrearsSinceLastPayment != 0 {
					t.Errorf("DaysInArrearsSinceLastPayment = %v, want %v", gotPeriod.DaysInArrearsSinceLastPayment, 0)
				}
				if !gotPeriod.DebtForArrearsSinceLastPayment.Equal(decimal.Zero) {
					t.Errorf("DebtForArrearsSinceLastPayment = %v, want %v", gotPeriod.DebtForArrearsSinceLastPayment, decimal.Zero)
				}
				if gotPeriod.TotalDaysInArrears != 0 {
					t.Errorf("TotalDaysInArrears = %v, want %v", gotPeriod.TotalDaysInArrears, 0)
				}
				if !gotPeriod.TotalDebtForArrears.Equal(decimal.Zero) {
					t.Errorf("TotalDebtForArrears = %v, want %v", gotPeriod.DebtForArrearsSinceLastPayment, decimal.Zero)
				}
				if !gotPeriod.TotalDebtOfPayment.Equal(got.PaymentAgreed) {
					t.Errorf("TotalDebtOfPayment = %v, want %v", gotPeriod.TotalDebtOfPayment, got.PaymentAgreed)
				}
				if !gotPeriod.TotalDebt.Equal(got.PaymentAgreed) {
					t.Errorf("TotalDebt = %v, want %v", gotPeriod.TotalDebt, got.PaymentAgreed)
				}
				if !gotPeriod.TotalPaid.Equal(decimal.Zero) {
					t.Errorf("TotalPaid = %v, want %v", gotPeriod.TotalPaid, decimal.Zero)
				}
				if !gotPeriod.TotalPaidToDebtForArrears.Equal(decimal.Zero) {
					t.Errorf("TotalPaidToDebtForArrears = %v, want %v", gotPeriod.TotalPaidToDebtForArrears, decimal.Zero)
				}
				if !gotPeriod.TotalPaidToRegularDebt.Equal(decimal.Zero) {
					t.Errorf("TotalPaidToRegularDebt = %v, want %v", gotPeriod.TotalPaidToRegularDebt, decimal.Zero)
				}
				if !gotPeriod.TotalPaidExtraToPrincipal.Equal(decimal.Zero) {
					t.Errorf("TotalPaidExtraToPrincipal = %v, want %v", gotPeriod.TotalPaidExtraToPrincipal, decimal.Zero)
				}
			}
		})
	}
}

func TestLoan_LiquidatePeriods(t *testing.T) {
	type liquidationExpected struct {
		periodNumber                   uint
		daysInArrearsSinceLastPayment  int
		debtForArrearsSinceLastPayment decimal.Decimal
		totalDebt                      decimal.Decimal
	}
	liqDate := toDate(2019, 3, 30)
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestNewLoanForCreate_LiquidatePeriods-1",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1, &liqDate},
		[]liquidationExpected{
			{1, 58, toDecimal(1498.6906), toDecimal(23646.8279)},
			{2, 30, toDecimal(775.1848), toDecimal(22923.3221)},
			{3, 0, decimal.Zero, toDecimal(22148.1373)},
			{4, 0, decimal.Zero, toDecimal(22148.1373)},
			{5, 0, decimal.Zero, toDecimal(22148.1373)},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoanForCreate(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, config.DefaultGraceDays, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoanForCreate() error = %v", err)
				return
			}
			got.LiquidateLoan(*tt.args.liquidationDate)
			for _, periodExpected := range tt.want {
				gotPeriod := got.periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if gotPeriod.DaysInArrearsSinceLastPayment != periodExpected.daysInArrearsSinceLastPayment {
					t.Errorf("DaysInArrearsSinceLastPayment = %v, want %v", gotPeriod.DaysInArrearsSinceLastPayment, periodExpected.daysInArrearsSinceLastPayment)
				}
				if !gotPeriod.DebtForArrearsSinceLastPayment.Equal(periodExpected.debtForArrearsSinceLastPayment) {
					t.Errorf("DebtForArrearsSinceLastPayment = %v, want %v", gotPeriod.DebtForArrearsSinceLastPayment, periodExpected.debtForArrearsSinceLastPayment)
				}
				if gotPeriod.TotalDaysInArrears != periodExpected.daysInArrearsSinceLastPayment {
					t.Errorf("TotalDaysInArrears = %v, want %v", gotPeriod.TotalDaysInArrears, periodExpected.daysInArrearsSinceLastPayment)
				}
				if !gotPeriod.TotalDebtForArrears.Equal(periodExpected.debtForArrearsSinceLastPayment) {
					t.Errorf("TotalDebtForArrears = %v, want %v", gotPeriod.TotalDebtForArrears, periodExpected.debtForArrearsSinceLastPayment)
				}
				if !gotPeriod.TotalDebt.Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt = %v, want %v", gotPeriod.TotalDebt, periodExpected.totalDebt)
				}
			}
		})
	}
}

func TestLoan_payLoan(t *testing.T) {
	type liquidationExpected struct {
		periodNumber                   uint
		daysInArrearsSinceLastPayment  int
		debtForArrearsSinceLastPayment decimal.Decimal
		totalDebt                      decimal.Decimal
	}
	liqDate := toDate(2019, 3, 30)
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestNewLoanForCreate_LiquidatePeriods-1",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1, &liqDate},
		[]liquidationExpected{
			{1, 58, toDecimal(1498.6906), toDecimal(23646.8279)},
			{2, 30, toDecimal(775.1848), toDecimal(22923.3221)},
			{3, 0, decimal.Zero, toDecimal(22148.1373)},
			{4, 0, decimal.Zero, toDecimal(22148.1373)},
			{5, 0, decimal.Zero, toDecimal(22148.1373)},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoanForCreate(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, config.DefaultGraceDays, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoanForCreate() error = %v", err)
				return
			}
			got.LiquidateLoan(*tt.args.liquidationDate)
			for _, periodExpected := range tt.want {
				gotPeriod := got.periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if gotPeriod.DaysInArrearsSinceLastPayment != periodExpected.daysInArrearsSinceLastPayment {
					t.Errorf("DaysInArrearsSinceLastPayment = %v, want %v", gotPeriod.DaysInArrearsSinceLastPayment, periodExpected.daysInArrearsSinceLastPayment)
				}
				if !gotPeriod.DebtForArrearsSinceLastPayment.Equal(periodExpected.debtForArrearsSinceLastPayment) {
					t.Errorf("DebtForArrearsSinceLastPayment = %v, want %v", gotPeriod.DebtForArrearsSinceLastPayment, periodExpected.debtForArrearsSinceLastPayment)
				}
				if gotPeriod.TotalDaysInArrears != periodExpected.daysInArrearsSinceLastPayment {
					t.Errorf("TotalDaysInArrears = %v, want %v", gotPeriod.TotalDaysInArrears, periodExpected.daysInArrearsSinceLastPayment)
				}
				if !gotPeriod.TotalDebtForArrears.Equal(periodExpected.debtForArrearsSinceLastPayment) {
					t.Errorf("TotalDebtForArrears = %v, want %v", gotPeriod.TotalDebtForArrears, periodExpected.debtForArrearsSinceLastPayment)
				}
				if !gotPeriod.TotalDebt.Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt = %v, want %v", gotPeriod.TotalDebt, periodExpected.totalDebt)
				}
			}
		})
	}
}

func toDecimal(number float64) decimal.Decimal {
	return decimal.NewFromFloat(number)
}

func toDate(year, month, day int) time.Time {
	return utils.DateWithoutTime(year, month, day)
}
