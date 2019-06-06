package application

import (
	clientDomain "loans/client/domain"
	loanDomain "loans/loan/domain"
	"loans/shared/config"
	"loans/shared/errors"
	"loans/shared/utils"
	"loans/shared/utils/financial"
	"time"

	"github.com/shopspring/decimal"
)

type LoanService struct {
	loanRepository   loanDomain.LoanRepository
	clientRepository clientDomain.ClientRepository
}

// NewLoanService creates a loan service with necessary dependencies.
func NewLoanService(loanRepository loanDomain.LoanRepository, clientRepository clientDomain.ClientRepository) LoanService {
	return LoanService{
		loanRepository:   loanRepository,
		clientRepository: clientRepository,
	}
}

func (s *LoanService) FindAllLoans() ([]loanDomain.Loan, error) {
	return s.loanRepository.FindAll()
}

func (s *LoanService) SimulateLoan(request CreateLoanRequest) (*LoanAmortizationsResponse, error) {
	loan, error := loanDomain.NewLoanForCreate(
		request.Principal,
		request.InterestRatePeriod,
		request.PeriodNumbers,
		request.StartDate,
		request.ClientID,
	)
	if error != nil {
		return nil, error
	}
	amortizations := financial.Amortizations(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers))
	var response LoanAmortizationsResponse
	response.ID = 0
	response.Principal = request.Principal
	response.InterestRatePeriod = request.InterestRatePeriod
	response.PeriodNumbers = request.PeriodNumbers
	response.StartDate = request.StartDate
	response.ClientID = request.ClientID
	response.PaymentAgreed = amortizations[0].Payment
	response.Amortizations = make([]AmortizationResponse, len(amortizations))
	for index, amoritzation := range amortizations {
		response.Amortizations[index].Period = index + 1
		response.Amortizations[index].PaymentDate = utils.AddMothToTimeForPayment(response.StartDate, index+1)
		response.Amortizations[index].InitialPrincipal = amoritzation.InitialPrincipal
		response.Amortizations[index].Payment = amoritzation.Payment
		response.Amortizations[index].InterestRatePeriod = amoritzation.InterestRatePeriod
		response.Amortizations[index].ToInterest = amoritzation.ToInterest
		response.Amortizations[index].ToPrincipal = amoritzation.ToPrincipal
		response.Amortizations[index].FinalPrincipal = amoritzation.FinalPrincipal
	}
	return &response, nil
}

func (s *LoanService) CreateLoan(request CreateLoanRequest) error {
	loan, error := loanDomain.NewLoanForCreate(
		request.Principal,
		request.InterestRatePeriod,
		request.PeriodNumbers,
		request.StartDate,
		request.ClientID,
	)
	if error != nil {
		return error
	}
	return s.loanRepository.StoreLoan(&loan)
}

func (s *LoanService) PayLoan(payment *loanDomain.Payment) error {
	loanPeriodsWithDue, error := s.loanRepository.FindBillsWithDueOrOpenOrderedByPeriodAsc(payment.LoanID)
	if error != nil {
		return error
	}
	if error := s.loanRepository.StorePayment(payment); error != nil {
		return error
	}
	remainingPayment := payment.PaymentAmount.RoundBank(config.Round)
	for index, period := range loanPeriodsWithDue {
		paymentToPeriod := decimal.Zero
		if remainingPayment.LessThanOrEqual(decimal.Zero) {
			break
		}
		if remainingPayment.LessThanOrEqual(period.TotalDue) || len(loanPeriodsWithDue) == (index+1) {
			paymentToPeriod = remainingPayment
		} else {
			paymentToPeriod = period.TotalDue
		}
		period.LiquidateBill(payment.PaymentDate)
		if error := s.payLoanPeriod(paymentToPeriod, period); error != nil {
			return error
		}
		remainingPayment = remainingPayment.Sub(paymentToPeriod).RoundBank(config.Round)
	}
	return nil
}

func (s *LoanService) closeLoan(loanID uint) error {
	loan, error := s.loanRepository.FindLoanByID(loanID)
	if error != nil {
		return error
	}
	loan.State = loanDomain.LoanStateClosed
	return s.loanRepository.UpdateLoan(&loan)
}

func (s *LoanService) initialLoanPeriod(loanID uint) error {
	loan, error := s.loanRepository.FindLoanByID(loanID)
	if error != nil {
		return error
	}
	bills, _ := s.loanRepository.FindBillsByLoanID(loanID)
	if len(bills) > 0 {
		return &errors.GracefulError{ErrorCode: errors.BillAlreadyExist}
	}
	period := 1
	newBill := loanDomain.LoanPeriod{}
	newBill.LoanID = loan.ID
	newBill.Period = uint(period)
	newBill.BillStartDate = loan.StartDate
	newBill.BillEndDate = utils.AddMothToTimeForPayment(newBill.BillStartDate, 1)
	balancePeriod := balanceExpectedInSpecificPeriodOfLoan(loan, period)
	fillDefaultAmountValues(&newBill, balancePeriod)
	return s.loanRepository.StoreBill(&newBill)
}

func (s *LoanService) recurringLoanPeriod(loanID uint) error {
	loan, error := s.loanRepository.FindLoanByID(loanID)
	if error != nil {
		return error
	}
	OldLoanPeriod := loanDomain.LoanPeriod{}
	OldLoanPeriod, error = s.loanRepository.FindBillOpenPeriodByLoanID(loanID)
	if error != nil {
		return error
	}
	if time.Now().Before(OldLoanPeriod.BillEndDate) {
		return nil
	}
	period := int(OldLoanPeriod.Period + 1)
	newBill := loanDomain.LoanPeriod{
		LoanID:        loan.ID,
		Period:        uint(period),
		BillStartDate: OldLoanPeriod.BillEndDate.AddDate(0, 0, 1),
		BillEndDate:   utils.AddMothToTimeForPayment(OldLoanPeriod.BillEndDate, 1),
	}
	nextBalance := nextBalanceFromLoanPeriod(OldLoanPeriod)
	fillDefaultAmountValues(&newBill, nextBalance)
	if error := s.loanRepository.StoreBill(&newBill); error != nil {
		return error
	}
	if error := s.closePeriod(&OldLoanPeriod); error != nil {
		return error
	}
	return s.recurringLoanPeriod(loanID)
}

func (s *LoanService) closePeriod(period *loanDomain.LoanPeriod) error {
	period.PeriodStatus = loanDomain.LoanPeriodStateClosed
	return s.loanRepository.UpdateBill(period)
}

func (s *LoanService) payLoanPeriod(periodPayment decimal.Decimal, loanPeriod loanDomain.LoanPeriod) error {
	loanPeriodMovement := new(loanDomain.LoanPeriodMovement)
	loanPeriodMovement.FillInitialBillMovementFromBill(loanPeriod)
	loanPeriod.ApplyPayment(periodPayment)
	if loanPeriod.FinalPrincipal.LessThanOrEqual(decimal.Zero) {
		loanPeriod.PeriodStatus = loanDomain.LoanPeriodStateClosed
		s.closeLoan(loanPeriod.LoanID)
	}
	loanPeriodMovement.FillFinalBillMovementFromBill(loanPeriod)
	if error := s.loanRepository.StoreBillMovement(loanPeriodMovement); error != nil {
		return error
	}
	return s.loanRepository.UpdateBill(&loanPeriod)
}

func balanceExpectedInSpecificPeriodOfLoan(loan loanDomain.Loan, period int) financial.Balance {
	return financial.BalanceExpectedInSpecificPeriod(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers), period)
}

func nextBalanceFromLoanPeriod(loanPeriod loanDomain.LoanPeriod) financial.Balance {
	balance := financial.Balance{}
	balance.InitialPrincipal = loanPeriod.InitialPrincipal
	balance.Payment = loanPeriod.Payment
	balance.InterestRatePeriod = loanPeriod.InterestRate
	balance.ToInterest = loanPeriod.InterestOfPayment
	balance.ToPrincipal = loanPeriod.PrincipalOfPayment
	balance.FinalPrincipal = loanPeriod.FinalPrincipal
	return financial.NextBalanceFromBefore(balance)
}

func fillDefaultAmountValues(bill *loanDomain.LoanPeriod, balance financial.Balance) {
	round := config.Round
	bill.State = loanDomain.LoanPeriodStateDue
	bill.PeriodStatus = loanDomain.PeriodStatusOpen
	bill.PaymentDate = bill.BillEndDate
	bill.InitialPrincipal = balance.InitialPrincipal
	bill.Payment = balance.Payment.RoundBank(round)
	bill.InterestOfPayment = balance.ToInterest.RoundBank(round)
	bill.InterestRate = balance.InterestRatePeriod.RoundBank(round)
	bill.PrincipalOfPayment = balance.ToPrincipal.RoundBank(round)
	bill.Paid = decimal.Zero
	bill.DaysLate = 0
	bill.FeeLateDue = decimal.Zero
	bill.PaymentDue = bill.Payment
	bill.TotalDue = bill.Payment
	bill.FinalPrincipal = balance.FinalPrincipal.RoundBank(round)
	bill.LastLiquidationDate = bill.PaymentDate
}
