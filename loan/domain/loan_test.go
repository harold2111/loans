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
	periodNumbers      int
	startDate          time.Time
	clientID           int
}

type liquidationExpected struct {
	periodNumber        int
	totalDaysInDefault  int
	TotalRegularDebt    decimal.Decimal
	totalDebtForDefault decimal.Decimal
	totalDebt           decimal.Decimal
	defaultExpecteds    int
	paymentExpecteds    int
	state               string
}

func TestLoan_CreateLoan(t *testing.T) {
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
			"TestLoan_CreateLoan-1",
			createLoanArgs{toDecimal(450000.0), toDecimal(0.05), 36, toDate(2019, 12, 16), 1},
			expected{toDecimal(27195.5057), toDate(2022, 12, 16)},
		},
		{
			"TestLoan_CreateLoan-2",
			createLoanArgs{toDecimal(1000.0), toDecimal(0.01), 12, toDate(2019, 12, 16), 1},
			expected{toDecimal(88.8488), toDate(2020, 12, 16)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
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
			periodsLen := len(got.Periods)
			if periodsLen != int(tt.args.periodNumbers) {
				t.Errorf("len(periods)  = %v, want %v", periodsLen, tt.args.periodNumbers)
			}
			lastPeriod := got.Periods[periodsLen-1]
			if !lastPeriod.FinalPrincipal.Equal(decimal.Zero) {
				t.Errorf("lastPeriod.FinalPrincipal  = %v, want %v", lastPeriod.FinalPrincipal, decimal.Zero)
			}
		})
	}
}

func TestLoan_CreatePeriods(t *testing.T) {
	type periodExpected struct {
		periodNumber       int
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
			"TestLoan_CreatePeriods-1",
			createLoanArgs{toDecimal(3500000), toDecimal(0.02), 5, toDate(2019, 1, 31), 1},
			[]periodExpected{
				{1, toDate(2019, 1, 31), toDate(2019, 2, 27), toDecimal(3500000), toDecimal(672554.3794), toDecimal(70000), toDecimal(2827445.6206)},
				{2, toDate(2019, 2, 28), toDate(2019, 3, 30), toDecimal(2827445.6206), toDecimal(686005.4670), toDecimal(56548.9124), toDecimal(2141440.1537)},
				{3, toDate(2019, 3, 31), toDate(2019, 4, 29), toDecimal(2141440.1537), toDecimal(699725.5763), toDecimal(42828.8031), toDecimal(1441714.5774)},
				{4, toDate(2019, 4, 30), toDate(2019, 5, 30), toDecimal(1441714.5774), toDecimal(713720.0878), toDecimal(28834.2915), toDecimal(727994.4896)},
				{5, toDate(2019, 5, 31), toDate(2019, 6, 29), toDecimal(727994.4896), toDecimal(727994.4896), toDecimal(14559.8898), toDecimal(0)},
			},
		},
		{
			"TestLoan_CreatePeriods-2",
			createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1},
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
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
				return
			}
			for _, periodExpected := range tt.want {
				gotPeriod := got.Periods[periodExpected.periodNumber-1]
				if gotPeriod.LoanID != got.ID {
					t.Errorf("LoanID = %v, want %v", gotPeriod.LoanID, got.ID)
				}
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if gotPeriod.State != PeriodStateOpen {
					t.Errorf("State = %v, want %v", gotPeriod.PeriodNumber, PeriodStateOpen)
				}
				if !gotPeriod.StartDate.Equal(periodExpected.startDate) {
					t.Errorf("StartDate = %v, want %v", gotPeriod.StartDate, periodExpected.startDate)
				}
				if !gotPeriod.EndDate.Equal(periodExpected.endDate) {
					t.Errorf("EndDate = %v, want %v", gotPeriod.EndDate, periodExpected.endDate)
				}
				maxPaymentDateExpected := periodExpected.endDate.AddDate(0, 0, config.DaysAfterEndDateToConsiderateInDefault)
				if !gotPeriod.MaxPaymentDate.Equal(maxPaymentDateExpected) {
					t.Errorf("MaxPaymentDate = %v, want %v", gotPeriod.MaxPaymentDate, maxPaymentDateExpected)
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
				//Modifiable fields
				if !gotPeriod.TotalPaidToRegularDebt.Equal(decimal.Zero) {
					t.Errorf("TotalPaidToRegularDebt = %v, want %v", gotPeriod.TotalPaidToRegularDebt, decimal.Zero)
				}
				if !gotPeriod.TotalPaidExtraToPrincipal.Equal(decimal.Zero) {
					t.Errorf("TotalPaidExtraToPrincipal = %v, want %v", gotPeriod.TotalPaidExtraToPrincipal, decimal.Zero)
				}
				//Calculated fields
				if !gotPeriod.lastLiquidationDate().Equal(periodExpected.endDate) {
					t.Errorf("lastLiquidationDate() = %v, want %v", gotPeriod.lastLiquidationDate(), periodExpected.endDate)
				}
				if !gotPeriod.TotalDefaultDebt().Equal(decimal.Zero) {
					t.Errorf("TotalDefaultDebt() = %v, want %v", gotPeriod.TotalDefaultDebt(), decimal.Zero)
				}
				if !gotPeriod.TotalRegularDebt().Equal(got.PaymentAgreed) {
					t.Errorf("TotalRegularDebt() = %v, want %v", gotPeriod.TotalRegularDebt(), got.PaymentAgreed)
				}
				if !gotPeriod.TotalDebt().Equal(got.PaymentAgreed) {
					t.Errorf("TotalDebt() = %v, want %v", gotPeriod.TotalDebt(), got.PaymentAgreed)
				}
				if gotPeriod.TotalDaysInDefault() != 0 {
					t.Errorf("TotalDaysInDefault() = %v, want %v", gotPeriod.TotalDaysInDefault(), 0)
				}

			}
		})
	}
}

func TestLoan_LiquidatePeriods(t *testing.T) {
	liqDate := toDate(2019, 3, 30)
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestLoan_LiquidatePeriods-1",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1},
		[]liquidationExpected{
			{1, 58, toDecimal(22148.1373), toDecimal(1498.6906), toDecimal(23646.8279), 1, 0, PeriodStateDue},
			{2, 30, toDecimal(22148.1373), toDecimal(775.1848), toDecimal(22923.3221), 1, 0, PeriodStateDue},
			{3, 0, toDecimal(22148.1373), decimal.Zero, toDecimal(22148.1373), 0, 0, PeriodStateDue},
			{4, 0, toDecimal(22148.1373), decimal.Zero, toDecimal(22148.1373), 0, 0, PeriodStateOpen},
			{5, 0, toDecimal(22148.1373), decimal.Zero, toDecimal(22148.1373), 0, 0, PeriodStateOpen},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
				return
			}
			got.LiquidateLoan(liqDate)
			for _, periodExpected := range tt.want {
				gotPeriod := got.Periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if !gotPeriod.TotalRegularDebt().Equal(periodExpected.TotalRegularDebt) {
					t.Errorf("TotalRegularDebt() = %v, want %v", gotPeriod.TotalRegularDebt(), periodExpected.TotalRegularDebt)
				}
				if gotPeriod.TotalDaysInDefault() != periodExpected.totalDaysInDefault {
					t.Errorf("TotalDaysInDefault() = %v, want %v", gotPeriod.TotalDaysInDefault(), periodExpected.totalDaysInDefault)
				}
				if !gotPeriod.TotalDefaultDebt().Equal(periodExpected.totalDebtForDefault) {
					t.Errorf("TotalDefaultDebt() = %v, want %v", gotPeriod.TotalDefaultDebt(), periodExpected.totalDebtForDefault)
				}
				if !gotPeriod.TotalDebt().Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt() = %v, want %v", gotPeriod.TotalDebt(), periodExpected.totalDebt)
				}
				if gotPeriod.State != periodExpected.state {
					t.Errorf("State = %v, want %v", gotPeriod.State, periodExpected.state)
				}
				if len(gotPeriod.DefaultPeriods) != periodExpected.defaultExpecteds {
					t.Errorf("len(gotPeriod.PeriodDefaults) = %v, want %v", len(gotPeriod.DefaultPeriods), periodExpected.defaultExpecteds)
				}
				if len(gotPeriod.Payments) != periodExpected.paymentExpecteds {
					t.Errorf("len(gotPeriod.PeriodPayments)  = %v, want %v", len(gotPeriod.Payments), periodExpected.paymentExpecteds)
				}
			}
		})
	}
}

func TestLoan_payLoanOnlyRegular(t *testing.T) {
	paymentDate := toDate(2019, 3, 30)
	payment := Payment{
		PaymentAmount:   toDecimal(48741.9365),
		RemainingAmount: toDecimal(48741.9365),
		PaymentDate:     paymentDate,
		PaymentType:     ExtraToNextPeriods,
	}
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestLoan_payLoanOnlyRegular",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1},
		[]liquidationExpected{
			{1, 58, decimal.Zero, decimal.Zero, decimal.Zero, 1, 1, PeriodStatePaid},
			{2, 30, decimal.Zero, decimal.Zero, decimal.Zero, 1, 1, PeriodStatePaid},
			{3, 0, toDecimal(19976.3508), decimal.Zero, toDecimal(19976.3508), 0, 1, PeriodStateDue},
			{4, 0, toDecimal(22148.1373), decimal.Zero, toDecimal(22148.1373), 0, 0, PeriodStateOpen},
			{5, 0, toDecimal(22148.1373), decimal.Zero, toDecimal(22148.1373), 0, 0, PeriodStateOpen},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
				return
			}
			got.ApplyPayment(payment)
			for _, periodExpected := range tt.want {
				gotPeriod := got.Periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if !gotPeriod.TotalRegularDebt().Equal(periodExpected.TotalRegularDebt) {
					t.Errorf("TotalRegularDebt() = %v, want %v", gotPeriod.TotalRegularDebt(), periodExpected.TotalRegularDebt)
				}
				if gotPeriod.TotalDaysInDefault() != periodExpected.totalDaysInDefault {
					t.Errorf("TotalDaysInDefault() = %v, want %v", gotPeriod.TotalDaysInDefault(), periodExpected.totalDaysInDefault)
				}
				if !gotPeriod.TotalDefaultDebt().Equal(periodExpected.totalDebtForDefault) {
					t.Errorf("TotalDefaultDebt() = %v, want %v", gotPeriod.TotalDefaultDebt(), periodExpected.totalDebtForDefault)
				}
				if !gotPeriod.TotalDebt().Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt() = %v, want %v", gotPeriod.TotalDebt(), periodExpected.totalDebt)
				}
				if gotPeriod.State != periodExpected.state {
					t.Errorf("State = %v, want %v", gotPeriod.State, periodExpected.state)
				}
				if len(gotPeriod.DefaultPeriods) != periodExpected.defaultExpecteds {
					t.Errorf("len(gotPeriod.PeriodDefaults) = %v, want %v", len(gotPeriod.DefaultPeriods), periodExpected.defaultExpecteds)
				}
				if len(gotPeriod.Payments) != periodExpected.paymentExpecteds {
					t.Errorf("len(gotPeriod.PeriodPayments)  = %v, want %v", len(gotPeriod.Payments), periodExpected.paymentExpecteds)
				}
			}
		})
	}
}

func TestLoan_payLoanInFirstMonthWithExtraToNextPeriods(t *testing.T) {
	paymentDate := toDate(2019, 1, 31)
	payment := Payment{
		PaymentAmount:   toDecimal(110740.6866),
		RemainingAmount: toDecimal(110740.6866),
		PaymentDate:     paymentDate,
		PaymentType:     ExtraToNextPeriods,
	}
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestLoan_payLoanOnlyRegular",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1},
		[]liquidationExpected{
			{1, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 1, PeriodStatePaid},
			{2, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 1, PeriodStatePaid},
			{3, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 1, PeriodStatePaid},
			{4, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 1, PeriodStatePaid},
			{5, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 1, PeriodStatePaid},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
				return
			}
			got.ApplyPayment(payment)
			for _, periodExpected := range tt.want {
				gotPeriod := got.Periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if !gotPeriod.TotalRegularDebt().Equal(periodExpected.TotalRegularDebt) {
					t.Errorf("TotalRegularDebt() = %v, want %v", gotPeriod.TotalRegularDebt(), periodExpected.TotalRegularDebt)
				}
				if gotPeriod.TotalDaysInDefault() != periodExpected.totalDaysInDefault {
					t.Errorf("TotalDaysInDefault() = %v, want %v", gotPeriod.TotalDaysInDefault(), periodExpected.totalDaysInDefault)
				}
				if !gotPeriod.TotalDefaultDebt().Equal(periodExpected.totalDebtForDefault) {
					t.Errorf("TotalDefaultDebt() = %v, want %v", gotPeriod.TotalDefaultDebt(), periodExpected.totalDebtForDefault)
				}
				if !gotPeriod.TotalDebt().Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt() = %v, want %v", gotPeriod.TotalDebt(), periodExpected.totalDebt)
				}
				if gotPeriod.State != periodExpected.state {
					t.Errorf("State = %v, want %v", gotPeriod.State, periodExpected.state)
				}
				if len(gotPeriod.DefaultPeriods) != periodExpected.defaultExpecteds {
					t.Errorf("len(gotPeriod.PeriodDefaults) = %v, want %v", len(gotPeriod.DefaultPeriods), periodExpected.defaultExpecteds)
				}
				if len(gotPeriod.Payments) != periodExpected.paymentExpecteds {
					t.Errorf("len(gotPeriod.PeriodPayments)  = %v, want %v", len(gotPeriod.Payments), periodExpected.paymentExpecteds)
				}
			}
		})
	}
}

func TestLoan_payLoanWithExraToPrincipalAnullingLastPeriod(t *testing.T) {
	paymentDate := toDate(2019, 3, 30)
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestLoan_payLoanWithExraToPrincipalAnullingLastPeriod-1",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1},
		[]liquidationExpected{
			{1, 58, decimal.Zero, decimal.Zero, decimal.Zero, 1, 1, PeriodStatePaid},
			{2, 30, decimal.Zero, decimal.Zero, decimal.Zero, 1, 1, PeriodStatePaid},
			{3, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 1, PeriodStatePaid},
			{4, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 2, PeriodStatePaid},
			{5, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 0, PeriodStateAnnulled},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
				return
			}
			payment := Payment{
				PaymentAmount:   toDecimal(132941.1144),
				RemainingAmount: toDecimal(132941.1144),
				PaymentDate:     paymentDate,
				PaymentType:     ExtraToPrincipal,
			}
			got.ApplyPayment(payment)
			for _, periodExpected := range tt.want {
				gotPeriod := got.Periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if !gotPeriod.TotalRegularDebt().Equal(periodExpected.TotalRegularDebt) {
					t.Errorf("TotalRegularDebt() = %v, want %v", gotPeriod.TotalRegularDebt(), periodExpected.TotalRegularDebt)
				}
				if gotPeriod.TotalDaysInDefault() != periodExpected.totalDaysInDefault {
					t.Errorf("TotalDaysInDefault() = %v, want %v", gotPeriod.TotalDaysInDefault(), periodExpected.totalDaysInDefault)
				}
				if !gotPeriod.TotalDefaultDebt().Equal(periodExpected.totalDebtForDefault) {
					t.Errorf("TotalDefaultDebt() = %v, want %v", gotPeriod.TotalDefaultDebt(), periodExpected.totalDebtForDefault)
				}
				if !gotPeriod.TotalDebt().Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt() = %v, want %v", gotPeriod.TotalDebt(), periodExpected.totalDebt)
				}
				if gotPeriod.State != periodExpected.state {
					t.Errorf("State = %v, want %v", gotPeriod.State, periodExpected.state)
				}
				if len(gotPeriod.DefaultPeriods) != periodExpected.defaultExpecteds {
					t.Errorf("len(gotPeriod.PeriodDefaults) = %v, want %v", len(gotPeriod.DefaultPeriods), periodExpected.defaultExpecteds)
				}
				if len(gotPeriod.Payments) != periodExpected.paymentExpecteds {
					t.Errorf("len(gotPeriod.PeriodPayments)  = %v, want %v", len(gotPeriod.Payments), periodExpected.paymentExpecteds)
				}
			}
		})
	}
}

func TestLoan_payWholePrincipalOnFirstMonth(t *testing.T) {
	paymentDate := toDate(2019, 1, 10)
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestLoan_payWholePrincipalOnFirstMonth-1",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1},
		[]liquidationExpected{
			{1, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 2, PeriodStatePaid},
			{2, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 0, PeriodStateAnnulled},
			{3, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 0, PeriodStateAnnulled},
			{4, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 0, PeriodStateAnnulled},
			{5, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 0, PeriodStateAnnulled},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
				return
			}
			payment := Payment{
				PaymentAmount:   toDecimal(103500.0),
				RemainingAmount: toDecimal(103500.0),
				PaymentDate:     paymentDate,
				PaymentType:     ExtraToPrincipal,
			}
			got.ApplyPayment(payment)
			for _, periodExpected := range tt.want {
				gotPeriod := got.Periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if !gotPeriod.TotalRegularDebt().Equal(periodExpected.TotalRegularDebt) {
					t.Errorf("TotalRegularDebt() = %v, want %v", gotPeriod.TotalRegularDebt(), periodExpected.TotalRegularDebt)
				}
				if gotPeriod.TotalDaysInDefault() != periodExpected.totalDaysInDefault {
					t.Errorf("TotalDaysInDefault() = %v, want %v", gotPeriod.TotalDaysInDefault(), periodExpected.totalDaysInDefault)
				}
				if !gotPeriod.TotalDefaultDebt().Equal(periodExpected.totalDebtForDefault) {
					t.Errorf("TotalDefaultDebt() = %v, want %v", gotPeriod.TotalDefaultDebt(), periodExpected.totalDebtForDefault)
				}
				if !gotPeriod.TotalDebt().Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt() = %v, want %v", gotPeriod.TotalDebt(), periodExpected.totalDebt)
				}
				if gotPeriod.State != periodExpected.state {
					t.Errorf("State = %v, want %v", gotPeriod.State, periodExpected.state)
				}
				if len(gotPeriod.DefaultPeriods) != periodExpected.defaultExpecteds {
					t.Errorf("len(gotPeriod.PeriodDefaults) = %v, want %v", len(gotPeriod.DefaultPeriods), periodExpected.defaultExpecteds)
				}
				if len(gotPeriod.Payments) != periodExpected.paymentExpecteds {
					t.Errorf("len(gotPeriod.PeriodPayments)  = %v, want %v", len(gotPeriod.Payments), periodExpected.paymentExpecteds)
				}
			}
		})
	}
}

func TestLoan_payInitialPrincipalOnFirstMonth(t *testing.T) {
	paymentDate := toDate(2019, 1, 10)
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestLoan_payInitialPrincipalOnFirstMonth-1",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1},
		[]liquidationExpected{
			{1, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 2, PeriodStatePaid},
			{2, 0, toDecimal(3622.500), decimal.Zero, toDecimal(3622.500), 0, 0, PeriodStateOpen},
			{3, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 0, PeriodStateAnnulled},
			{4, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 0, PeriodStateAnnulled},
			{5, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 0, PeriodStateAnnulled},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
				return
			}
			payment := Payment{
				PaymentAmount:   toDecimal(100000.0),
				RemainingAmount: toDecimal(100000.0),
				PaymentDate:     paymentDate,
				PaymentType:     ExtraToPrincipal,
			}
			got.ApplyPayment(payment)
			for _, periodExpected := range tt.want {
				gotPeriod := got.Periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if !gotPeriod.TotalRegularDebt().Equal(periodExpected.TotalRegularDebt) {
					t.Errorf("TotalRegularDebt() = %v, want %v", gotPeriod.TotalRegularDebt(), periodExpected.TotalRegularDebt)
				}
				if gotPeriod.TotalDaysInDefault() != periodExpected.totalDaysInDefault {
					t.Errorf("TotalDaysInDefault() = %v, want %v", gotPeriod.TotalDaysInDefault(), periodExpected.totalDaysInDefault)
				}
				if !gotPeriod.TotalDefaultDebt().Equal(periodExpected.totalDebtForDefault) {
					t.Errorf("TotalDefaultDebt() = %v, want %v", gotPeriod.TotalDefaultDebt(), periodExpected.totalDebtForDefault)
				}
				if !gotPeriod.TotalDebt().Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt() = %v, want %v", gotPeriod.TotalDebt(), periodExpected.totalDebt)
				}
				if gotPeriod.State != periodExpected.state {
					t.Errorf("State = %v, want %v", gotPeriod.State, periodExpected.state)
				}
				if len(gotPeriod.DefaultPeriods) != periodExpected.defaultExpecteds {
					t.Errorf("len(gotPeriod.PeriodDefaults) = %v, want %v", len(gotPeriod.DefaultPeriods), periodExpected.defaultExpecteds)
				}
				if len(gotPeriod.Payments) != periodExpected.paymentExpecteds {
					t.Errorf("len(gotPeriod.PeriodPayments)  = %v, want %v", len(gotPeriod.Payments), periodExpected.paymentExpecteds)
				}
			}
		})
	}
}

func TestLoan_payLoanPartialPaymentInThirdMonth(t *testing.T) {
	paymentDate1 := toDate(2019, 3, 30)
	payment1 := Payment{
		PaymentAmount:   toDecimal(48741.9365),
		RemainingAmount: toDecimal(48741.9365),
		PaymentDate:     paymentDate1,
		PaymentType:     ExtraToNextPeriods,
	}
	paymentDate2 := toDate(2019, 3, 30)
	payment2 := Payment{
		PaymentAmount:   toDecimal(10000),
		RemainingAmount: toDecimal(10000),
		PaymentDate:     paymentDate2,
		PaymentType:     ExtraToNextPeriods,
	}
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestLoan_payLoanOnlyRegular",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1},
		[]liquidationExpected{
			{1, 58, decimal.Zero, decimal.Zero, decimal.Zero, 1, 1, PeriodStatePaid},
			{2, 30, decimal.Zero, decimal.Zero, decimal.Zero, 1, 1, PeriodStatePaid},
			{3, 0, toDecimal(9976.3508), decimal.Zero, toDecimal(9976.3508), 0, 2, PeriodStateDue},
			{4, 0, toDecimal(22148.1373), decimal.Zero, toDecimal(22148.1373), 0, 0, PeriodStateOpen},
			{5, 0, toDecimal(22148.1373), decimal.Zero, toDecimal(22148.1373), 0, 0, PeriodStateOpen},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
				return
			}
			got.ApplyPayment(payment1)
			got.ApplyPayment(payment2)
			for _, periodExpected := range tt.want {
				gotPeriod := got.Periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if !gotPeriod.TotalRegularDebt().Equal(periodExpected.TotalRegularDebt) {
					t.Errorf("TotalRegularDebt() = %v, want %v", gotPeriod.TotalRegularDebt(), periodExpected.TotalRegularDebt)
				}
				if gotPeriod.TotalDaysInDefault() != periodExpected.totalDaysInDefault {
					t.Errorf("TotalDaysInDefault() = %v, want %v", gotPeriod.TotalDaysInDefault(), periodExpected.totalDaysInDefault)
				}
				if !gotPeriod.TotalDefaultDebt().Equal(periodExpected.totalDebtForDefault) {
					t.Errorf("TotalDefaultDebt() = %v, want %v", gotPeriod.TotalDefaultDebt(), periodExpected.totalDebtForDefault)
				}
				if !gotPeriod.TotalDebt().Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt() = %v, want %v", gotPeriod.TotalDebt(), periodExpected.totalDebt)
				}
				if gotPeriod.State != periodExpected.state {
					t.Errorf("State = %v, want %v", gotPeriod.State, periodExpected.state)
				}
				if len(gotPeriod.DefaultPeriods) != periodExpected.defaultExpecteds {
					t.Errorf("len(gotPeriod.PeriodDefaults) = %v, want %v", len(gotPeriod.DefaultPeriods), periodExpected.defaultExpecteds)
				}
				if len(gotPeriod.Payments) != periodExpected.paymentExpecteds {
					t.Errorf("len(gotPeriod.PeriodPayments)  = %v, want %v", len(gotPeriod.Payments), periodExpected.paymentExpecteds)
				}
			}
		})
	}
}

func TestLoan_payAllInThirdMothWithSecondPayment(t *testing.T) {
	paymentDate1 := toDate(2019, 3, 30)
	payment1 := Payment{
		PaymentAmount:   toDecimal(48741.9365),
		RemainingAmount: toDecimal(48741.9365),
		PaymentDate:     paymentDate1,
		PaymentType:     ExtraToNextPeriods,
	}
	paymentDate2 := toDate(2019, 3, 30)
	payment2 := Payment{
		PaymentAmount:   toDecimal(64272.6254),
		RemainingAmount: toDecimal(64272.6254),
		PaymentDate:     paymentDate2,
		PaymentType:     ExtraToNextPeriods,
	}
	tests := []struct {
		name string
		args createLoanArgs
		want []liquidationExpected
	}{{
		"TestLoan_payLoanOnlyRegular",
		createLoanArgs{toDecimal(100000), toDecimal(0.035), 5, toDate(2019, 1, 1), 1},
		[]liquidationExpected{
			{1, 58, decimal.Zero, decimal.Zero, decimal.Zero, 1, 1, PeriodStatePaid},
			{2, 30, decimal.Zero, decimal.Zero, decimal.Zero, 1, 1, PeriodStatePaid},
			{3, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 2, PeriodStatePaid},
			{4, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 1, PeriodStatePaid},
			{5, 0, decimal.Zero, decimal.Zero, decimal.Zero, 0, 1, PeriodStatePaid},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoan(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoan() error = %v", err)
				return
			}
			got.ApplyPayment(payment1)
			got.ApplyPayment(payment2)
			for _, periodExpected := range tt.want {
				gotPeriod := got.Periods[periodExpected.periodNumber-1]
				if gotPeriod.PeriodNumber != periodExpected.periodNumber {
					t.Errorf("PeriodNumber = %v, want %v", gotPeriod.PeriodNumber, periodExpected.periodNumber)
				}
				if !gotPeriod.TotalRegularDebt().Equal(periodExpected.TotalRegularDebt) {
					t.Errorf("TotalRegularDebt() = %v, want %v", gotPeriod.TotalRegularDebt(), periodExpected.TotalRegularDebt)
				}
				if gotPeriod.TotalDaysInDefault() != periodExpected.totalDaysInDefault {
					t.Errorf("TotalDaysInDefault() = %v, want %v", gotPeriod.TotalDaysInDefault(), periodExpected.totalDaysInDefault)
				}
				if !gotPeriod.TotalDefaultDebt().Equal(periodExpected.totalDebtForDefault) {
					t.Errorf("TotalDefaultDebt() = %v, want %v", gotPeriod.TotalDefaultDebt(), periodExpected.totalDebtForDefault)
				}
				if !gotPeriod.TotalDebt().Equal(periodExpected.totalDebt) {
					t.Errorf("TotalDebt() = %v, want %v", gotPeriod.TotalDebt(), periodExpected.totalDebt)
				}
				if gotPeriod.State != periodExpected.state {
					t.Errorf("State = %v, want %v", gotPeriod.State, periodExpected.state)
				}
				if len(gotPeriod.DefaultPeriods) != periodExpected.defaultExpecteds {
					t.Errorf("len(gotPeriod.PeriodDefaults) = %v, want %v", len(gotPeriod.DefaultPeriods), periodExpected.defaultExpecteds)
				}
				if len(gotPeriod.Payments) != periodExpected.paymentExpecteds {
					t.Errorf("len(gotPeriod.PeriodPayments)  = %v, want %v", len(gotPeriod.Payments), periodExpected.paymentExpecteds)
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
