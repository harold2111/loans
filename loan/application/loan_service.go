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
func (s *LoanService) FindAllLoans() ([]LoanResponse, error) {
	var responses []LoanResponse
	loans, error := s.loanRepository.FindAll()
	if error != nil {
		return responses, nil
	}
	for _, loan := range loans {
		response := transformLoanToLoanResponse(loan)
		responses = append(responses, response)
	}
	return responses, nil
}

//SimulateLoan creates and returns a loan without persisting
func (s *LoanService) SimulateLoan(request CreateLoanRequest) (LoanResponse, error) {
	var response LoanResponse
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
	response = transformLoanToLoanResponse(loan)
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
	payment, error := loanDomain.NewPayment(
		request.LoanID,
		request.PaymentAmount,
		request.PaymentDate,
		request.PaymentType,
	)
	if error != nil {
		return PayLoanResponse{}, error
	}
	loanRepository := s.loanRepository
	loan, error := loanRepository.FindLoanByID(payment.LoanID)
	if error != nil {
		return response, error
	}
	error = s.loanRepository.StorePayment(&payment)
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
	response.PaymentDate = remainingPayment.PaymentDate
	response.PaymentAmount = remainingPayment.PaymentAmount
	response.RemainingAmount = remainingPayment.RemainingAmount
	return response, nil
}

func transformLoanToLoanResponse(loan loanDomain.Loan) LoanResponse {
	var response LoanResponse
	response.ID = 0
	response.Principal = loan.Principal
	response.InterestRatePeriod = loan.InterestRatePeriod
	response.PeriodNumbers = loan.PeriodNumbers
	response.StartDate = loan.StartDate
	response.ClientID = loan.ClientID
	response.PaymentAgreed = loan.PaymentAgreed
	periods := loan.Periods
	response.Periods = make([]PeriodResponse, len(periods))
	for index, period := range periods {
		response.Periods[index].Period = index + 1
		response.Periods[index].MaxPaymentDate = period.MaxPaymentDate
		response.Periods[index].InitialPrincipal = period.InitialPrincipal
		response.Periods[index].Payment = period.Payment
		response.Periods[index].InterestRatePeriod = period.InterestRate
		response.Periods[index].ToInterest = period.InterestOfPayment
		response.Periods[index].ToPrincipal = period.PrincipalOfPayment
		response.Periods[index].FinalPrincipal = period.FinalPrincipal
	}
	return response
}
