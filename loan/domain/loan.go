package domain

import (
	"sort"
	"time"

	"github.com/harold2111/loans/shared/config"
	"github.com/harold2111/loans/shared/utils"
	"github.com/harold2111/loans/shared/utils/financial"

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
	GraceDays          uint
	ClientID           uint `gorm:"not null"`
}

func NewLoanForCreate(
	principal decimal.Decimal,
	interestRatePeriod decimal.Decimal,
	periodNumbers uint,
	startDate time.Time,
	graceDays uint,
	clientID uint) (Loan, error) {

	loan := Loan{
		Principal:          principal,
		InterestRatePeriod: interestRatePeriod,
		PeriodNumbers:      periodNumbers,
		StartDate:          startDate,
		GraceDays:          graceDays,
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

func (l *Loan) LiquidateLoan(liquidationDate time.Time) {
	for periodIndex := 0; periodIndex < len(l.periods); periodIndex++ {
		l.periods[periodIndex].LiquidateByDate(liquidationDate, l.GraceDays)
	}
}

func (l *Loan) ApplyPayment(payment Payment) decimal.Decimal {
	periods := l.periods
	numPeriods := len(periods)
	l.LiquidateLoan(payment.PaymentDate)
	sort.Slice(periods, func(p, q int) bool { return periods[p].PeriodNumber < periods[q].PeriodNumber })
	remainingPayment := payment.PaymentAmount
	//REGULAR PAYMENT
	if remainingPayment.GreaterThan(decimal.Zero) {
		for periodIndex := 0; periodIndex < numPeriods; periodIndex++ {
			period := &periods[periodIndex]
			if period.State == LoanPeriodStateDue ||
				(period.State == LoanPeriodStateOpen && payment.PaymentType == ExtraToNextPeriods) {
				remainingPayment = period.ApplyPayment(payment.ID, remainingPayment)
			}
		}
	}
	//EXTRA PRINCIPAL PAYMENT
	firstOpenPeriod := l.findFirstOpenPeriod()
	if remainingPayment.GreaterThan(decimal.Zero) && firstOpenPeriod != nil && payment.PaymentType == ExtraToPrincipal {
		remainingPayment = firstOpenPeriod.ApplyPayment(payment.ID, remainingPayment)
		if remainingPayment.GreaterThanOrEqual(decimal.Zero) {
			remainingPayment = firstOpenPeriod.ApplyPaymentToPrincipal(payment.ID, remainingPayment)
			l.recalculatePeriodsForExtraPrincipalPayment(*firstOpenPeriod)
		}
	}
	l.roundDecimalValues()
	return remainingPayment
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
		var endDate time.Time
		periodNumber := index + 1
		if index == 0 {
			startDate = l.StartDate
		} else {
			startDate = utils.AddMothToTimeForPayment(l.StartDate, periodNumber-1)
		}
		endDate = utils.AddMothToTimeForPayment(l.StartDate, periodNumber).AddDate(0, 0, -1)
		periods[index].PeriodNumber = uint(periodNumber)
		periods[index].State = LoanPeriodStateOpen
		periods[index].StartDate = startDate
		periods[index].EndDate = endDate
		periods[index].PaymentDate = endDate
		periods[index].InitialPrincipal = amoritzation.InitialPrincipal
		periods[index].Payment = amoritzation.Payment
		periods[index].InterestRate = amoritzation.InterestRatePeriod
		periods[index].PrincipalOfPayment = amoritzation.ToPrincipal
		periods[index].InterestOfPayment = amoritzation.ToInterest
		periods[index].FinalPrincipal = amoritzation.FinalPrincipal
		periods[index].LastPaymentDate = endDate
		periods[index].TotalDebtOfPayment = amoritzation.Payment
		periods[index].TotalDebt = amoritzation.Payment

	}
	l.periods = periods
}

func (l *Loan) recalculatePeriodsForExtraPrincipalPayment(periodWithExtraPrincialPayment LoanPeriod) {
	periods := l.periods
	numPeriods := len(periods)
	recalculatedPeriodIndex := int(periodWithExtraPrincialPayment.PeriodNumber)
	beforePeriodIndex := recalculatedPeriodIndex - 1
	anullateRestOfPeriods := false
	for recalculatedPeriodIndex < numPeriods {
		beforePeriod := &periods[beforePeriodIndex]
		recalculatedPeriod := &periods[recalculatedPeriodIndex]
		if beforePeriod.FinalPrincipal.LessThanOrEqual(decimal.Zero) && !anullateRestOfPeriods {
			anullateRestOfPeriods = true
		}
		if anullateRestOfPeriods {
			recalculatedPeriod.State = LoanPeriodStateAnnuelled
		} else {
			recalculatedPeriod.InitialPrincipal = beforePeriod.FinalPrincipal
			recalculatedPeriod.InterestOfPayment = recalculatedPeriod.InitialPrincipal.Mul(recalculatedPeriod.InterestRate)
			recalculatedPeriodTotalPayment := recalculatedPeriod.InitialPrincipal.Add(recalculatedPeriod.InterestOfPayment)
			if recalculatedPeriodTotalPayment.LessThanOrEqual(recalculatedPeriod.Payment) {
				recalculatedPeriod.Payment = recalculatedPeriodTotalPayment
			}
			recalculatedPeriod.PrincipalOfPayment = recalculatedPeriod.Payment.Sub(recalculatedPeriod.InterestOfPayment)
			recalculatedPeriod.FinalPrincipal = recalculatedPeriod.InitialPrincipal.Sub(recalculatedPeriod.PrincipalOfPayment)
		}
		beforePeriodIndex++
		recalculatedPeriodIndex++
	}
}

func (l *Loan) findFirstOpenPeriod() *LoanPeriod {
	numPeriods := len(l.periods)
	periods := l.periods
	for index := 0; index < numPeriods; index++ {
		period := &periods[index]
		if period.State == LoanPeriodStateOpen {
			return period
		}
	}
	return nil
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
		periods[index].TotalDebt = periods[index].TotalDebt.RoundBank(config.Round)

	}
	l.periods = periods
}
