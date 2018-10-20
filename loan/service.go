package loan

import (
	"loans/loan/dtos"
	"loans/models"
)

// LoanService is the interface that provides loan user case methods.
type LoanService interface {
	SimulateLoan(request dtos.CreateLoanRequest) (*dtos.LoanAmortizationsResponse, error)
	FindAllLoans() ([]models.Loan, error)
	CreateLoan(loan *models.Loan) error
	PayLoan(payment *models.Payment) error
}
