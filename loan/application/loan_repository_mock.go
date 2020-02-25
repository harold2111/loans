package application

import "github.com/harold2111/loans/loan/domain"

type LoanRepositoryMock struct {
	StoreLoanMock func(loan *domain.Loan) error
}

func (r *LoanRepositoryMock) FindAll() ([]domain.Loan, error) {
	panic("not implemented") // TODO: Implement
}

func (r *LoanRepositoryMock) StoreLoan(loan *domain.Loan) error {
	return r.StoreLoanMock(loan)
}

func (r *LoanRepositoryMock) UpdateLoan(loan *domain.Loan) error {
	panic("not implemented") // TODO: Implement
}

func (r *LoanRepositoryMock) FindLoanByID(loanID uint) (domain.Loan, error) {
	panic("not implemented") // TODO: Implement
}

func (r *LoanRepositoryMock) StoreBill(bill *domain.Period) error {
	panic("not implemented") // TODO: Implement
}

func (r *LoanRepositoryMock) UpdateBill(bill *domain.Period) error {
	panic("not implemented") // TODO: Implement
}

func (r *LoanRepositoryMock) FindBillsByLoanID(loanID uint) ([]domain.Period, error) {
	panic("not implemented") // TODO: Implement
}

func (r *LoanRepositoryMock) FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID uint) ([]domain.Period, error) {
	panic("not implemented") // TODO: Implement
}

func (r *LoanRepositoryMock) FindBillOpenPeriodByLoanID(loanID uint) (domain.Period, error) {
	panic("not implemented") // TODO: Implement
}

func (r *LoanRepositoryMock) StorePayment(payment *domain.Payment) error {
	panic("not implemented") // TODO: Implement
}
