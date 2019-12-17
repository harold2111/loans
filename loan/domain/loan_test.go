package domain

import (
	"testing"
	"time"

	"github.com/harold2111/loans/shared/config"
	"github.com/shopspring/decimal"
)

func TestNewLoanForCreate(t *testing.T) {
	type args struct {
		principal          decimal.Decimal
		interestRatePeriod decimal.Decimal
		periodNumbers      uint
		startDate          time.Time
		clientID           uint
	}
	type expected struct {
		principal          decimal.Decimal
		interestRatePeriod decimal.Decimal
		periodNumbers      uint
		startDate          time.Time
		paymentAgreed      decimal.Decimal
		state              string
		closeDateAgreed    time.Time
		clientID           uint
	}
	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			"CreateLoanTest-1",
			args{decimal.NewFromFloat(450000.0), decimal.NewFromFloat(0.05), 36, parseDateWithoutTime("2019-12-16"), 1},
			expected{decimal.NewFromFloat(450000.0), decimal.NewFromFloat(0.05), 36, parseDateWithoutTime("2019-12-16"), decimal.NewFromFloat(27195.5057), LoanStateActive, parseDateWithoutTime("2022-12-16"), 1},
		},
		{
			"CreateLoanTest-2",
			args{decimal.NewFromFloat(1000.0), decimal.NewFromFloat(0.01), 12, parseDateWithoutTime("2019-12-16"), 1},
			expected{decimal.NewFromFloat(1000.0), decimal.NewFromFloat(0.01), 12, parseDateWithoutTime("2019-12-16"), decimal.NewFromFloat(88.8488), LoanStateActive, parseDateWithoutTime("2020-12-16"), 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLoanForCreate(tt.args.principal, tt.args.interestRatePeriod, tt.args.periodNumbers, tt.args.startDate, tt.args.clientID)
			if err != nil {
				t.Errorf("NewLoanForCreate() error = %v", err)
				return
			}
			if !got.Principal.Equal(tt.want.principal) {
				t.Errorf("Principal = %v, want %v", got.Principal, tt.want.principal)
			}
			if !got.InterestRatePeriod.Equal(tt.want.interestRatePeriod) {
				t.Errorf("InterestRatePeriod = %v, want %v", got.InterestRatePeriod, tt.want.interestRatePeriod)
			}
			if got.PeriodNumbers != tt.want.periodNumbers {
				t.Errorf("PeriodNumbers = %v, want %v", got.PeriodNumbers, tt.want.periodNumbers)
			}
			if !got.PaymentAgreed.Equal(tt.want.paymentAgreed) {
				t.Errorf("PaymentAgreed = %v, want %v", got.PaymentAgreed, tt.want.paymentAgreed)
			}
			if !got.StartDate.Equal(tt.want.startDate) {
				t.Errorf("StartDate = %v, want %v", got.StartDate, tt.want.startDate)
			}
			if got.State != tt.want.state {
				t.Errorf("State = %v, want %v", got.StartDate, tt.want.startDate)
			}
			if !got.CloseDateAgreed.Equal(tt.want.closeDateAgreed) {
				t.Errorf("CloseDateAgreed = %v, want %v", got.CloseDateAgreed, tt.want.closeDateAgreed)
			}
			periodsLen := len(got.periods)
			if periodsLen != int(tt.want.periodNumbers) {
				t.Errorf("len(periods)  = %v, want %v", periodsLen, tt.want.periodNumbers)
			}
			lastPeriod := got.periods[periodsLen-1]
			if !lastPeriod.FinalPrincipal.Equal(decimal.Zero) {
				t.Errorf("lastPeriod.FinalPrincipal  = %v, want %v", lastPeriod.FinalPrincipal, decimal.Zero)
			}
		})
	}
}

func parseDateWithoutTime(timeString string) time.Time {
	const RFC3339FullDate = "2006-01-02"
	parsedTime, _ := time.ParseInLocation(RFC3339FullDate, timeString, config.DefaultLocation())
	return parsedTime
}

func TestLoan_validateForCreation(t *testing.T) {
	tests := []struct {
		name    string
		l       *Loan
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.validateForCreation(); (err != nil) != tt.wantErr {
				t.Errorf("Loan.validateForCreation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoan_calculatePaymentAgreed(t *testing.T) {
	tests := []struct {
		name string
		l    *Loan
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.calculatePaymentAgreed()
		})
	}
}

func TestLoan_calculateCloseDateAgreed(t *testing.T) {
	tests := []struct {
		name string
		l    *Loan
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.calculateCloseDateAgreed()
		})
	}
}

func TestLoan_calculatePeriods(t *testing.T) {
	tests := []struct {
		name string
		l    *Loan
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.calculatePeriods()
		})
	}
}
