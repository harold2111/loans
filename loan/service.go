package loan

import (
	"loans/financial"
	"loans/models"
)

// LoanService is the interface that provides loan user case methods.
type LoanService interface {
	SimulateLoan(loan models.Loan) []financial.Balance
	FindAllLoans() ([]models.Loan, error)
	CreateLoan(loan *models.Loan) error
	PayLoan(payment *models.Payment) error
}
