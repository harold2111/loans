package loan

// Repository provides access a loan store.
type Repository interface {
	StoreLoan(loan *Loan) error
	UpdateLoan(loan *Loan) error
	FindLoanByID(loanID uint) (Loan, error)
	/************/
	StoreBill(bill *Bill) error
	UpdateBill(bill *Bill) error
	FindBillsByLoanID(loanID uint) ([]Bill, error)
	FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID uint) ([]Bill, error)
	FindBillOpenPeriodByLoanID(loanID uint) (Bill, error)
	StoreBillMovement(billMovement *BillMovement) error
	/************/
	StorePayment(payment *Payment) error
}
