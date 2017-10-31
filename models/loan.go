package models

import (
	"loans/config"
	"loans/errors"
	"loans/financial"
	"loans/utils"
	"time"

	"github.com/shopspring/decimal"

	"github.com/jinzhu/gorm"
)

const (
	LoanStateActive = "ACTIVE"
)

type Loan struct {
	gorm.Model
	Principal          decimal.Decimal `gorm:"type:numeric"`
	InterestRatePeriod decimal.Decimal `gorm:"type:numeric"`
	PeriodNumbers      uint
	PaymentAgreed      decimal.Decimal `gorm:"type:numeric"`
	StartDate          time.Time
	CloseDateAgreed    time.Time
	CloseDate          *time.Time
	State              string
	Client             Client
	ClientID           uint `gorm:"not null"`
}

func (loan *Loan) Create() error {
	loan.State = LoanStateActive
	loan.StartDate = loan.StartDate.In(config.DefaultLocation())
	calculatePaymentOfLoan(loan)
	calculateCloseDateAgreed(loan)
	if error := config.DB.Create(loan).Error; error != nil {
		return error
	}
	if error := CreateInitialBill(loan.ID); error != nil {
		return nil
	}
	if error := RecurringLoanBillingByLoanID(loan.ID); error != nil {
		return nil
	}
	return nil
}

func FindLoanByID(loanID uint) (Loan, error) {
	var loan Loan
	response := config.DB.First(&loan, loanID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{loanID}
			return loan, &errors.RecordNotFound{ErrorCode: errors.ClientNotExist, MessagesParameters: messagesParameters}
		}
		return loan, error
	}
	return loan, nil
}

func calculateCloseDateAgreed(loan *Loan) {
	loan.CloseDateAgreed = utils.AddMothToTimeForPayment(loan.StartDate, int(loan.PeriodNumbers))
}

func calculatePaymentOfLoan(loan *Loan) {
	loan.PaymentAgreed = financial.CalculatePayment(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers)).RoundBank(config.Round)

}

func balanceExpectedInSpecificPeriodOfLoan(loan Loan, period int) financial.Balance {
	return financial.BalanceExpectedInSpecificPeriod(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers), period)
}
