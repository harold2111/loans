package domain

import (
	"loans/shared/config"
	"loans/shared/utils"
	"loans/shared/utils/financial"
	"time"

	"github.com/shopspring/decimal"
)

const (
	LoanStateActive = "ACTIVE"
	LoanStateClosed = "CLOSED"
)

type Loan struct {
	ID                 uint `gorm:"primary_key"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time      `sql:"index"`
	Principal          decimal.Decimal `gorm:"type:numeric"`
	InterestRatePeriod decimal.Decimal `gorm:"type:numeric"`
	PeriodNumbers      uint
	PaymentAgreed      decimal.Decimal `gorm:"type:numeric"`
	StartDate          time.Time
	CloseDateAgreed    time.Time
	CloseDate          *time.Time
	State              string
	ClientID           uint `gorm:"not null"`
}

func NewLoanForCreate(
	principal decimal.Decimal,
	interestRatePeriod decimal.Decimal,
	periodNumbers uint,
	startDate time.Time,
	clientID uint) (Loan, error) {

	loan := Loan{
		Principal:          principal,
		InterestRatePeriod: interestRatePeriod,
		PeriodNumbers:      periodNumbers,
		StartDate:          startDate,
		ClientID:           clientID,
	}
	if error := loan.validateForCreation(); error != nil {
		return Loan{}, error
	}
	loan.StartDate = loan.StartDate.In(config.DefaultLocation())
	loan.calculatePaymentAgreed()
	loan.calculateCloseDateAgreed()
	loan.State = LoanStateActive
	return loan, nil
}

func (l *Loan) validateForCreation() error {
	if error := utils.ValidateVar("principal", l.Principal, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("interestRatePeriod", l.InterestRatePeriod, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("periodNumbers", l.PeriodNumbers, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("startDate", l.StartDate, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("clientID", l.ClientID, "required"); error != nil {
		return error
	}
	return nil
}

func (l *Loan) calculatePaymentAgreed() {
	l.PaymentAgreed = financial.CalculatePayment(l.Principal, l.InterestRatePeriod, int(l.PeriodNumbers)).RoundBank(config.Round)
}

func (l *Loan) calculateCloseDateAgreed() {
	l.CloseDateAgreed = utils.AddMothToTimeForPayment(l.StartDate, int(l.PeriodNumbers))
}
