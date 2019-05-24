package postgres

import (
	loanDomain "loans/loan/domain"
	"loans/shared/errors"
	"loans/shared/models"

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

func (r *loanRepository) FindAll() ([]models.Loan, error) {
	var loans []models.Loan
	response := r.db.Find(&loans)
	if error := response.Error; error != nil {
		return nil, error
	}
	return loans, nil
}

func (r *loanRepository) FindLoanByID(loanID uint) (models.Loan, error) {
	var loan models.Loan
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

func (r *loanRepository) StoreLoan(loan *models.Loan) error {
	return r.db.Create(loan).Error
}

func (r *loanRepository) UpdateLoan(loan *models.Loan) error {
	return r.db.Save(loan).Error
}

func (r *loanRepository) StoreBill(bill *models.Bill) error {
	return r.db.Create(bill).Error
}

func (r *loanRepository) UpdateBill(bill *models.Bill) error {
	return r.db.Save(bill).Error
}

func (r *loanRepository) FindBillsByLoanID(loanID uint) ([]models.Bill, error) {
	var bills []models.Bill
	r.db.Find(&bills, "loan_id = ?", loanID)
	return bills, nil
}

func (r *loanRepository) FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID uint) ([]models.Bill, error) {
	var bills []models.Bill
	r.db.Order("period").Find(&bills, "loan_id = ? AND state = ? OR period_status = ?", loanID, models.BillStateDue, models.PeriodStatusOpen)
	return bills, nil
}

func (r *loanRepository) FindBillOpenPeriodByLoanID(loanID uint) (models.Bill, error) {
	bill := models.Bill{}
	error := r.db.Raw("SELECT * FROM bills WHERE loan_id = ? AND period_status = ? AND period = (SELECT max(period) FROM bills where loan_id = ?)",
		loanID, models.PeriodStatusOpen, loanID).Scan(&bill).Error
	return bill, error
}

func (r *loanRepository) StoreBillMovement(billMovement *models.BillMovement) error {
	return r.db.Create(billMovement).Error
}

func (r *loanRepository) StorePayment(payment *models.Payment) error {
	return r.db.Create(payment).Error
}
