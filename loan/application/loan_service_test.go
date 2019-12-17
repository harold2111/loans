package application

import (
	"testing"
	"time"

	"github.com/harold2111/loans/loan/domain"
	"github.com/shopspring/decimal"
)

func TestCreateLoan(t *testing.T) {
	r := &LoanRepositoryMock{
		StoreLoanMock: func(loan *domain.Loan) error {
			return nil
		},
	}
	s := &clientRepositoryMock{}
	service := NewLoanService(r, s)
	createLoanRequest := CreateLoanRequest{
		Principal:          decimal.NewFromFloat(1000.0),
		InterestRatePeriod: decimal.NewFromFloat(1.0),
		PeriodNumbers:      12,
		StartDate:          time.Now(),
		ClientID:           1,
	}
	error := service.CreateLoan(createLoanRequest)
	if error != nil {
		t.Fatalf("CreateLoan throw exception %v", error)
	}
}
