package domain

// LoanRepository interface that provides access a loan store.
type LoanRepository interface {
	FindAll() ([]Loan, error)
	StoreLoan(loan *Loan) error
	UpdateLoan(loan *Loan) error
	FindLoanByID(loanID int) (Loan, error)
	StoreBill(bill *Period) error
	UpdateBill(bill *Period) error
	FindBillsByLoanID(loanID int) ([]Period, error)
	FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID int) ([]Period, error)
	FindBillOpenPeriodByLoanID(loanID int) (Period, error)
	StorePayment(payment *Payment) error
}
