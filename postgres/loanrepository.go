package postgres

import (
	"loans/errors"
	"loans/loan"

	"github.com/jinzhu/gorm"
)

type loanRepository struct {
	db *gorm.DB
}

// NewLoanRepository returns a new instance of a Postgres Loan repository.
func NewLoanRepository(db *gorm.DB) (loan.Repository, error) {
	r := &loanRepository{
		db: db,
	}
	return r, nil
}

func (r *loanRepository) StoreLoan(loan *loan.Loan) error {
	return r.db.Create(loan).Error
}

func (r *loanRepository) UpdateLoan(loan *loan.Loan) error {
	return r.db.Save(loan).Error
}

func (r *loanRepository) FindLoanByID(loanID uint) (loan.Loan, error) {
	var loan loan.Loan
	respose := r.db.First(&loan, loanID)
	if error := respose.Error; error != nil {
		if respose.RecordNotFound() {
			messagesParameters := []interface{}{loanID}
			return loan, &errors.RecordNotFound{ErrorCode: errors.ClientNotExist, MessagesParameters: messagesParameters}
		}
		return loan, error
	}
	return loan, nil
}

func (r *loanRepository) StoreBill(bill *loan.Bill) error {
	return r.db.Create(bill).Error
}

func (r *loanRepository) UpdateBill(bill *loan.Bill) error {
	return r.db.Save(bill).Error
}

func (r *loanRepository) FindBillsByLoanID(loanID uint) ([]loan.Bill, error) {
	var bills []loan.Bill
	r.db.Find(&bills, "loan_id = ?", loanID)
	return bills, nil
}

func (r *loanRepository) FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID uint) ([]loan.Bill, error) {
	var bills []loan.Bill
	r.db.Order("period").Find(&bills, "loan_id = ? AND state = ? OR period_status = ?", loanID, loan.BillStateDue, loan.PeriodStatusOpen)
	return bills, nil
}

func (r *loanRepository) FindBillOpenPeriodByLoanID(loanID uint) (loan.Bill, error) {
	bill := loan.Bill{}
	error := r.db.Raw("SELECT * FROM bills WHERE loan_id = ? AND period_status = ? AND period = (SELECT max(period) FROM bills where loan_id = ?)",
		loanID, loan.PeriodStatusOpen, loanID).Scan(&bill).Error
	return bill, error
}

func (r *loanRepository) StoreBillMovement(billMovement *loan.BillMovement) error {
	return r.db.Create(billMovement).Error
}

func (r *loanRepository) StorePayment(payment *loan.Payment) error {
	return r.db.Create(payment).Error
}
