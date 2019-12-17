package application

import (
	clientDomain "github.com/harold2111/loans/client/domain"
	loanDomain "github.com/harold2111/loans/loan/domain"
	"github.com/harold2111/loans/shared/config"
	"github.com/harold2111/loans/shared/utils"
	"github.com/harold2111/loans/shared/utils/financial"

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
	loanRepository := s.loanRepository
	loanPeriodsWithDebt, error := loanRepository.FindBillsWithDueOrOpenOrderedByPeriodAsc(payment.LoanID)
	if error != nil {
		return error
	}
	//Liquidate Periods
	/*
	Como nota primero deberia aplicar los pagos corrientes y luego deberia decidir que hacer con el excedente
	en caso de que se haya pagado algun extra.
	*/
	for _, period := range loanPeriodsWithDebt {
		period.LiquidateByDate(payment.PaymentDate)
	}
	var paidPeriods []loanDomain.LoanPeriod
	var paidPeriodMovements []loanDomain.LoanPeriodMovement
	remainingPayment := payment.PaymentAmount.RoundBank(config.Round)
	continueApplyingPayment := true
	for _, period := range loanPeriodsWithDebt {
		paymentToPeriod := decimal.Zero
		if continueApplyingPayment {
			if remainingPayment.LessThanOrEqual(period.TotalDebt) {
				paymentToPeriod = remainingPayment
			} else {
				paymentToPeriod = period.TotalDebt
			}
			loanPeriodMovement := period.ApplyPayment(payment.ID, paymentToPeriod)
			paidPeriods = append(paidPeriods, period)                             //I will use it for a store in batch
			paidPeriodMovements = append(paidPeriodMovements, loanPeriodMovement) //I will use it for a store in batch
			if period.FinalPrincipal.LessThanOrEqual(decimal.Zero) {
				s.closeLoan(period.LoanID)
			}
			remainingPayment = remainingPayment.Sub(paymentToPeriod).RoundBank(config.Round)
			continueApplyingPayment = remainingPayment.LessThanOrEqual(decimal.Zero)
		}
	}

	/*if error := loanRepository.StorePayment(payment); error != nil {
		return error
	}
	if error := s.loanRepository.StoreBillMovement(&loanPeriodMovement); error != nil {
		return error
	}
	if error := s.loanRepository.UpdateBill(&period); error != nil {
		return error
	}*/
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
