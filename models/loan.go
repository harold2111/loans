package models

import (
	"fmt"
	"loans/config"
	"loans/errors"
	"loans/financial"
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
	//TODO: validate clientID
	calculatePaymentOfLoan(loan)
	calculateCloseDateAgreed(loan)
	loan.State = LoanStateActive
	if error := config.DB.Create(loan).Error; error != nil {
		return nil
	}
	if error := CreateInitialBill(loan.ID); error != nil {
		return nil
	}
	if error := RecurringLoanBillingByLoanID(loan.ID); error != nil {
		return nil
	}
	if error := RecurringLoanBillingByLoanID(loan.ID); error != nil {
		return nil
	}
	return nil
}

func calculatePaymentOfLoan(loan *Loan) {
	loan.PaymentAgreed = financial.CalculatePayment(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers)).RoundBank(config.Round)

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
	fmt.Println("********************loan.StartDate: ", loan.StartDate)
	loan.CloseDateAgreed = addMothToTimeUtil(loan.StartDate, int(loan.PeriodNumbers))
	fmt.Println("********************loan.StartDate: ", loan.StartDate)
}

func addMothToTimeUtil(startTime time.Time, monthToAdd int) time.Time {
	endTime := startTime.AddDate(0, monthToAdd, 0)
	endTimeWithLastMothDay := time.Date(endTime.Year(), endTime.Month(), 0,
		endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Nanosecond(), endTime.Location())
	if startTime.Day() > endTimeWithLastMothDay.Day() {
		return endTimeWithLastMothDay
	}
	return endTime
}

func balanceExpectedInSpecificPeriodOfLoan(loan Loan, period int) financial.Balance {
	return financial.BalanceExpectedInSpecificPeriod(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers), period)
}
