package loan

import "loans/models"

// LoanService is the interface that provides loan user case methods.
type LoanService interface {
	CreateLoan(loan *models.Loan) error
	PayLoan(payment *models.Payment) error
}
