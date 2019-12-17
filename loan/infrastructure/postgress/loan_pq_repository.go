package postgres

import (
	loanDomain "github.com/harold2111/loans/loan/domain"
	"github.com/harold2111/loans/shared/errors"

	"github.com/jinzhu/gorm"
)

type loanRepository struct {
	db *gorm.DB
}

// NewLoanRepository returns a new instance of a Postgres Loan repository.
func NewLoanRepository(db *gorm.DB) (loanDomain.LoanRepository, error) {
	r := &loanRepository{
		db: db,
	}
	return r, nil
}

func (r *loanRepository) FindAll() ([]loanDomain.Loan, error) {
	var loans []loanDomain.Loan
	response := r.db.Find(&loans)
	if error := response.Error; error != nil {
		return nil, error
	}
	return loans, nil
}

func (r *loanRepository) FindLoanByID(loanID uint) (loanDomain.Loan, error) {
	var loan loanDomain.Loan
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

func (r *loanRepository) StoreLoan(loan *loanDomain.Loan) error {
	return r.db.Create(loan).Error
}

func (r *loanRepository) UpdateLoan(loan *loanDomain.Loan) error {
	return r.db.Save(loan).Error
}

func (r *loanRepository) StoreBill(bill *loanDomain.LoanPeriod) error {
	return r.db.Create(bill).Error
}

func (r *loanRepository) UpdateBill(bill *loanDomain.LoanPeriod) error {
	return r.db.Save(bill).Error
}

func (r *loanRepository) FindBillsByLoanID(loanID uint) ([]loanDomain.LoanPeriod, error) {
	var bills []loanDomain.LoanPeriod
	r.db.Find(&bills, "loan_id = ?", loanID)
	return bills, nil
}

func (r *loanRepository) FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID uint) ([]loanDomain.LoanPeriod, error) {
	var bills []loanDomain.LoanPeriod
	r.db.Order("period").Find(&bills, "loan_id = ? AND state = ? OR period_status = ?", loanID, loanDomain.LoanPeriodStateDue, loanDomain.LoanPeriodStateOpen)
	return bills, nil
}

func (r *loanRepository) FindBillOpenPeriodByLoanID(loanID uint) (loanDomain.LoanPeriod, error) {
	bill := loanDomain.LoanPeriod{}
	error := r.db.Raw("SELECT * FROM bills WHERE loan_id = ? AND period_status = ? AND period = (SELECT max(period) FROM bills where loan_id = ?)",
		loanID, loanDomain.LoanPeriodStateOpen, loanID).Scan(&bill).Error
	return bill, error
}

func (r *loanRepository) StoreBillMovement(billMovement *loanDomain.LoanPeriodMovement) error {
	return r.db.Create(billMovement).Error
}

func (r *loanRepository) StorePayment(payment *loanDomain.Payment) error {
	return r.db.Create(payment).Error
}
