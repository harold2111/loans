package domain

// LoanRepository interface that provides access a loan store.
type LoanRepository interface {
	FindAll() ([]Loan, error)
	StoreLoan(loan *Loan) error
	UpdateLoan(loan *Loan) error
	FindLoanByID(loanID uint) (Loan, error)
	StoreBill(bill *Period) error
	UpdateBill(bill *Period) error
	FindBillsByLoanID(loanID uint) ([]Period, error)
	FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID uint) ([]Period, error)
	FindBillOpenPeriodByLoanID(loanID uint) (Period, error)
	StorePayment(payment *Payment) error
}
