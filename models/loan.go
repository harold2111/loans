package models

import (
	"loans/config"
	"loans/errors"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"github.com/jinzhu/gorm"
)

const (
	LoanStateActive = "ACTIVE"
)

type Loan struct {
	gorm.Model
	Principal       decimal.Decimal `gorm:"type:numeric"`
	InterestRate    decimal.Decimal `gorm:"type:numeric"`
	PeriodNumbers   uint
	PaymentAgreed   decimal.Decimal `gorm:"type:numeric"`
	PaymentDate     time.Time
	StartDate       time.Time
	CloseDateAgreed time.Time
	CloseDate       *time.Time
	State           string
	Client          Client
	ClientID        uint `gorm:"not null"`
}

func (loan *Loan) Save() error {
	//TODO: validate clientID
	calculatePaymentOfLoan(loan)
	calculateCloseDateAgreed(loan)
	loan.State = LoanStateActive
	error := config.DB.Create(loan).Error
	return error
}

func calculatePaymentOfLoan(loan *Loan) {
	loan.PaymentAgreed = calculatePayment(loan.Principal, loan.InterestRate, int(loan.PeriodNumbers))

}

func FindLoanByID(loanID uint) (*Loan, error) {
	var client Loan
	response := config.DB.First(&client, loanID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{loanID}
			return nil, &errors.RecordNotFound{ErrorCode: errors.ClientNotExist, MessagesParameters: messagesParameters}
		}
		return nil, error
	}
	return &client, nil
}

func calculateCloseDateAgreed(loan *Loan) {
	loan.CloseDateAgreed = addMothToTimeUtil(loan.StartDate, int(loan.PeriodNumbers))
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

/*********************/
func calculatePayment(principal decimal.Decimal, interestRate decimal.Decimal, periodNumbers int) decimal.Decimal {

	hundred := decimal.NewFromFloat(100)
	one := decimal.NewFromFloat(1)
	n, _ := decimal.NewFromString(strconv.Itoa(periodNumbers))
	nNeg := n.Neg()
	rate := interestRate.Div(hundred)

	rateMulPrincipal := rate.Mul(principal)
	ratePlusOne := rate.Add(one)
	ratePlusOnePowNNeg := ratePlusOne.Pow(nNeg)
	oneMinusRatePlusOnePowNNeg := one.Sub(ratePlusOnePowNNeg)

	payment := rateMulPrincipal.Div(oneMinusRatePlusOnePowNNeg)

	return payment.RoundBank(5)
}

func getBalancingInSpecificPeriodNumber(principal decimal.Decimal, interestRate decimal.Decimal, periodNumbers int, specificPeriod int) decimal.Decimal {
	payment := calculatePayment(principal, interestRate, periodNumbers)
	hundred := decimal.NewFromFloat(100)
	interestDecimal := interestRate.Div(hundred)
	initialBalance := principal
	var finalBalance decimal.Decimal
	for period := 0; period < specificPeriod; period++ {
		toInterest := initialBalance.Mul(interestDecimal)
		toCapital := payment.Sub(toInterest)
		finalBalance = initialBalance.Sub(toCapital)
		initialBalance = finalBalance
	}
	return finalBalance.RoundBank(5)
}
