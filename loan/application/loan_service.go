package application

import (
	clientDomain "github.com/harold2111/loans/client/domain"
	loanDomain "github.com/harold2111/loans/loan/domain"
)

//LoanService service to operate over loans
type LoanService struct {
	loanRepository   loanDomain.LoanRepository
	clientRepository clientDomain.ClientRepository
}

// NewLoanService creates a loan service with its necessary dependencies.
func NewLoanService(loanRepository loanDomain.LoanRepository, clientRepository clientDomain.ClientRepository) LoanService {
	return LoanService{
		loanRepository:   loanRepository,
		clientRepository: clientRepository,
	}
}

//FindAllLoans finds all loans in the system
func (s *LoanService) FindAllLoans() ([]loanDomain.Loan, error) {
	return s.loanRepository.FindAll()
}

//SimulateLoan creates and returns a loan without persisting
func (s *LoanService) SimulateLoan(request CreateLoanRequest) (LoanAmortizationsResponse, error) {
	var response LoanAmortizationsResponse
	loan, error := loanDomain.NewLoan(
		request.Principal,
		request.InterestRatePeriod,
		request.PeriodNumbers,
		request.StartDate,
		request.ClientID,
	)
	if error != nil {
		return response, error
	}
	response.ID = 0
	response.Principal = loan.Principal
	response.InterestRatePeriod = loan.InterestRatePeriod
	response.PeriodNumbers = loan.PeriodNumbers
	response.StartDate = loan.StartDate
	response.ClientID = loan.ClientID
	response.PaymentAgreed = loan.PaymentAgreed
	periods := loan.Periods
	response.Amortizations = make([]AmortizationResponse, len(periods))
	for index, period := range periods {
		response.Amortizations[index].Period = index + 1
		response.Amortizations[index].MaxPaymentDate = period.MaxPaymentDate
		response.Amortizations[index].InitialPrincipal = period.InitialPrincipal
		response.Amortizations[index].Payment = period.Payment
		response.Amortizations[index].InterestRatePeriod = period.InterestRate
		response.Amortizations[index].ToInterest = period.InterestOfPayment
		response.Amortizations[index].ToPrincipal = period.PrincipalOfPayment
		response.Amortizations[index].FinalPrincipal = period.FinalPrincipal
	}
	return response, nil
}

//CreateLoan create and persist a new loan
func (s *LoanService) CreateLoan(request CreateLoanRequest) error {
	loan, error := loanDomain.NewLoan(
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

//PayLoan receive a payment for a specific loan
func (s *LoanService) PayLoan(request PayLoanRequest) (PayLoanResponse, error) {
	var response PayLoanResponse
	payment := loanDomain.NewPayment(
		request.LoanID,
		request.PaymentAmount,
		request.PaymentDate,
		request.PaymentType,
	)
	loanRepository := s.loanRepository
	loan, error := loanRepository.FindLoanByID(payment.LoanID)
	if error != nil {
		return response, error
	}
	remainingPayment := loan.ApplyPayment(payment)

	error = s.loanRepository.UpdateLoan(&loan)
	if error != nil {
		return response, error
	}
	response.ID = remainingPayment.ID
	response.LoanID = remainingPayment.LoanID
	response.PaymentAmount = remainingPayment.PaymentAmount
	response.RemainingAmount = remainingPayment.RemainingAmount
	return response, nil
}
