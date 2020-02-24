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
	//LoanStateActive represents a active loan.
	LoanStateActive = "ACTIVE"
	//LoanStateClosed represents a active closed.
	LoanStateClosed = "CLOSED"
)

//Loan represents a loan.
type Loan struct {
	ID                 int `gorm:"primary_key"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time      `sql:"index"`
	Principal          decimal.Decimal `gorm:"type:numeric"`
	InterestRatePeriod decimal.Decimal `gorm:"type:numeric"`
	PeriodNumbers      int
	PaymentAgreed      decimal.Decimal `gorm:"type:numeric"`
	StartDate          time.Time
	CloseDateAgreed    time.Time
	CloseDate          *time.Time
	State              string
	periods            []LoanPeriod
	ClientID           int `gorm:"not null"`
}

// NewLoanForCreate create a new Loan.
func NewLoanForCreate(
	principal decimal.Decimal,
	interestRatePeriod decimal.Decimal,
	periodNumbers int,
	startDate time.Time,
	clientID int) (Loan, error) {

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

// LiquidateLoan liquidate the loan periods for a specific date.
func (l *Loan) LiquidateLoan(liquidationDate time.Time) {
	for periodIndex := 0; periodIndex < len(l.periods); periodIndex++ {
		l.periods[periodIndex].liquidateByDate(liquidationDate)
	}
}

// ApplyPayment apply a payment to loan periods.
func (l *Loan) ApplyPayment(payment Payment) Payment {
	remainingPayment := payment
	periods := l.periods
	l.LiquidateLoan(payment.PaymentDate)
	sort.Slice(periods, func(p, q int) bool { return periods[p].PeriodNumber < periods[q].PeriodNumber })
	remainingPayment = l.applyRegularPayments(payment)
	remainingPayment = l.applyExtraToPrincipalPayments(payment)
	l.roundDecimalValues()
	return remainingPayment
}

func (l *Loan) applyRegularPayments(payment Payment) Payment {
	if payment.PaymentAmount.LessThanOrEqual(decimal.Zero) {
		return payment
	}
	remainingPayment := payment
	periods := l.periods
	numPeriods := len(periods)
	for periodIndex := 0; periodIndex < numPeriods; periodIndex++ {
		period := &periods[periodIndex]
		if period.State == LoanPeriodStateDue ||
			(period.State == LoanPeriodStateOpen && payment.isExtraToNextPeriods()) {
			remainingPayment = period.applyRegularPayment(remainingPayment)
		}
	}
	return remainingPayment
}

func (l *Loan) applyExtraToPrincipalPayments(payment Payment) Payment {
	remainingPayment := payment
	firstOpenPeriod := l.findFirstOpenPeriod()
	if remainingPayment.PaymentAmount.LessThanOrEqual(decimal.Zero) || !payment.isExtraToPrincipal() || firstOpenPeriod == nil {
		return remainingPayment
	}
	remainingPayment = firstOpenPeriod.applyRegularPayment(remainingPayment)
	if remainingPayment.PaymentAmount.GreaterThanOrEqual(decimal.Zero) {
		remainingPayment = firstOpenPeriod.applyPaymentToPrincipal(payment.ID, remainingPayment)
		l.recalculatePeriodsForExtraPrincipalPayment(*firstOpenPeriod)
	}
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
	for index, amortization := range amortizations {
		var startDate time.Time
		var endDate time.Time
		periodNumber := index + 1
		if index == 0 {
			startDate = l.StartDate
		} else {
			startDate = utils.AddMothToTimeForPayment(l.StartDate, periodNumber-1)
		}
		endDate = utils.AddMothToTimeForPayment(l.StartDate, periodNumber).AddDate(0, 0, -1)
		maxPaymentDate := endDate.AddDate(0, 0, config.DaysAfterEndDateToConsiderateInArrears)
		periods[index].PeriodNumber = periodNumber
		periods[index].State = LoanPeriodStateOpen
		periods[index].StartDate = startDate
		periods[index].EndDate = endDate
		periods[index].MaxPaymentDate = maxPaymentDate
		periods[index].InitialPrincipal = amortization.InitialPrincipal
		periods[index].Payment = amortization.Payment
		periods[index].InterestRate = amortization.InterestRatePeriod
		periods[index].PrincipalOfPayment = amortization.ToPrincipal
		periods[index].InterestOfPayment = amortization.ToInterest
		periods[index].FinalPrincipal = amortization.FinalPrincipal

	}
	l.periods = periods
}

func (l *Loan) recalculatePeriodsForExtraPrincipalPayment(periodWithExtraPrincipalPayment LoanPeriod) {
	periods := l.periods
	numPeriods := len(periods)
	recalculatedPeriodIndex := int(periodWithExtraPrincipalPayment.PeriodNumber)
	beforePeriodIndex := recalculatedPeriodIndex - 1
	annulateRestOfPeriods := false
	for recalculatedPeriodIndex < numPeriods {
		beforePeriod := &periods[beforePeriodIndex]
		recalculatedPeriod := &periods[recalculatedPeriodIndex]
		if beforePeriod.FinalPrincipal.LessThanOrEqual(decimal.Zero) && !annulateRestOfPeriods {
			annulateRestOfPeriods = true
		}
		if annulateRestOfPeriods {
			recalculatedPeriod.State = LoanPeriodStateAnnulled
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
	for index, amortization := range periods {
		periods[index].InitialPrincipal = periods[index].InitialPrincipal.RoundBank(config.Round)
		periods[index].Payment = periods[index].Payment.RoundBank(config.Round)
		periods[index].InterestRate = periods[index].InterestRate.RoundBank(config.Round)
		periods[index].PrincipalOfPayment = periods[index].PrincipalOfPayment.RoundBank(config.Round)
		periods[index].InterestOfPayment = periods[index].InterestOfPayment.RoundBank(config.Round)
		periods[index].FinalPrincipal = amortization.FinalPrincipal.RoundBank(config.Round)

	}
	l.periods = periods
}
