package domain

// LoanRepository interface that provides access a loan store.
type LoanRepository interface {
	FindAll() ([]Loan, error)
	StoreLoan(loan *Loan) error
	UpdateLoan(loan *Loan) error
	FindLoanByID(loanID uint) (Loan, error)
	StoreBill(bill *LoanPeriod) error
	UpdateBill(bill *LoanPeriod) error
	FindBillsByLoanID(loanID uint) ([]LoanPeriod, error)
	FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID uint) ([]LoanPeriod, error)
	FindBillOpenPeriodByLoanID(loanID uint) (LoanPeriod, error)
	StoreBillMovement(billMovement *LoanPeriodMovement) error
	StorePayment(payment *Payment) error
}
