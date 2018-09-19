package loan

import "loans/models"

// Repository provides access a loan store.
type LoanRepository interface {
	StoreLoan(loan *models.Loan) error
	UpdateLoan(loan *models.Loan) error
	FindLoanByID(loanID uint) (models.Loan, error)
	StoreBill(bill *models.Bill) error
	UpdateBill(bill *models.Bill) error
	FindBillsByLoanID(loanID uint) ([]models.Bill, error)
	FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID uint) ([]models.Bill, error)
	FindBillOpenPeriodByLoanID(loanID uint) (models.Bill, error)
	StoreBillMovement(billMovement *models.BillMovement) error
	StorePayment(payment *models.Payment) error
}
