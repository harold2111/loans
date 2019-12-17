package domain

import (
	"github.com/harold2111/loans/shared/config"
	"github.com/harold2111/loans/shared/utils"
	"github.com/harold2111/loans/shared/utils/financial"
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
	periods            []LoanPeriod
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
	loan.calculatePeriods()
	loan.roundDecimalValues()
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
	l.PaymentAgreed = financial.CalculatePayment(l.Principal, l.InterestRatePeriod, int(l.PeriodNumbers))
}

func (l *Loan) calculateCloseDateAgreed() {
	l.CloseDateAgreed = utils.AddMothToTimeForPayment(l.StartDate, int(l.PeriodNumbers))
}

func (l *Loan) calculatePeriods() {
	amortizations := financial.Amortizations(l.Principal, l.InterestRatePeriod, int(l.PeriodNumbers))
	periods := make([]LoanPeriod, len(amortizations))
	for index, amoritzation := range amortizations {
		var startDate time.Time
		if index == 0 {
			startDate = l.StartDate
		} else {
			startDate = utils.AddMothToTimeForPayment(l.StartDate, index+1)
		}
		endDate := utils.AddMothToTimeForPayment(startDate, index+1)
		periods[index].PeriodNumber = uint(index) + 1
		periods[index].State = LoanPeriodStateDue
		periods[index].StartDate = startDate
		periods[index].EndDate = endDate
		periods[index].PaymentDate = endDate
		periods[index].InitialPrincipal = amoritzation.InitialPrincipal
		periods[index].Payment = amoritzation.Payment
		periods[index].InterestRate = amoritzation.InterestRatePeriod
		periods[index].PrincipalOfPayment = amoritzation.ToPrincipal
		periods[index].InterestOfPayment = amoritzation.ToInterest
		periods[index].FinalPrincipal = amoritzation.FinalPrincipal
		periods[index].LastLiquidationDate = endDate
		periods[index].TotalDebtOfPayment = amoritzation.Payment

	}
	l.periods = periods
}

func (l *Loan) roundDecimalValues() {
	l.PaymentAgreed = l.PaymentAgreed.RoundBank(config.Round)
	periods := l.periods
	for index, amoritzation := range periods {
		periods[index].InitialPrincipal = periods[index].InitialPrincipal.RoundBank(config.Round)
		periods[index].Payment = periods[index].Payment.RoundBank(config.Round)
		periods[index].InterestRate = periods[index].InterestRate.RoundBank(config.Round)
		periods[index].PrincipalOfPayment = periods[index].PrincipalOfPayment.RoundBank(config.Round)
		periods[index].InterestOfPayment = periods[index].InterestOfPayment.RoundBank(config.Round)
		periods[index].FinalPrincipal = amoritzation.FinalPrincipal.RoundBank(config.Round)
		periods[index].TotalDebtOfPayment = periods[index].TotalDebtOfPayment.RoundBank(config.Round)

	}
	l.periods = periods
}

func (l *Loan) LiquidateLoan(liquidationDate time.Time) {
	for _, period := range l.periods {
		period.LiquidateByDate(liquidationDate)
	}
}
