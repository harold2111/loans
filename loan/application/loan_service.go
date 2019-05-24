package application

import (
	clientDomain "loans/client/domain"
	loanDomain "loans/loan/domain"
	"loans/shared/config"
	"loans/shared/errors"
	"loans/shared/models"
	"loans/shared/utils"
	"loans/shared/utils/financial"
	"time"

	"github.com/jinzhu/copier"
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

func (s *LoanService) FindAllLoans() ([]models.Loan, error) {
	return s.loanRepository.FindAll()
}

func (s *LoanService) SimulateLoan(request CreateLoanRequest) (*LoanAmortizationsResponse, error) {
	if error := utils.ValidateStruct(request); error != nil {
		return nil, error
	}
	loan := models.Loan{}
	if error := copier.Copy(&loan, &request); error != nil {
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

func (s *LoanService) CreateLoan(loan *models.Loan) error {
	loan.State = models.LoanStateActive
	loan.StartDate = loan.StartDate.In(config.DefaultLocation())
	calculatePaymentOfLoan(loan)
	calculateCloseDateAgreed(loan)
	if error := s.loanRepository.StoreLoan(loan); error != nil {
		return error
	}
	if error := s.initialBill(loan.ID); error != nil {
		return nil
	}
	if error := s.recurringBill(loan.ID); error != nil {
		return nil
	}
	return nil
}

func (s *LoanService) closeLoan(loanID uint) error {
	loan, error := s.loanRepository.FindLoanByID(loanID)
	if error != nil {
		return error
	}
	loan.State = models.LoanStateClosed
	return s.loanRepository.UpdateLoan(&loan)
}

func (s *LoanService) initialBill(loanID uint) error {
	loan, error := s.loanRepository.FindLoanByID(loanID)
	if error != nil {
		return error
	}
	bills, _ := s.loanRepository.FindBillsByLoanID(loanID)
	if len(bills) > 0 {
		return &errors.GracefulError{ErrorCode: errors.BillAlreadyExist}
	}
	period := 1
	newBill := models.Bill{}
	newBill.LoanID = loan.ID
	newBill.Period = uint(period)
	newBill.BillStartDate = loan.StartDate
	newBill.BillEndDate = utils.AddMothToTimeForPayment(newBill.BillStartDate, 1)
	balancePeriod := balanceExpectedInSpecificPeriodOfLoan(loan, period)
	fillDefaultAmountValues(&newBill, balancePeriod)
	return s.loanRepository.StoreBill(&newBill)
}

func (s *LoanService) recurringBill(loanID uint) error {
	loan, error := s.loanRepository.FindLoanByID(loanID)
	if error != nil {
		return error
	}
	oldBill := models.Bill{}
	oldBill, error = s.loanRepository.FindBillOpenPeriodByLoanID(loanID)
	if error != nil {
		return error
	}
	if time.Now().Before(oldBill.BillEndDate) {
		return nil
	}
	period := int(oldBill.Period + 1)
	newBill := models.Bill{
		LoanID:        loan.ID,
		Period:        uint(period),
		BillStartDate: oldBill.BillEndDate.AddDate(0, 0, 1),
		BillEndDate:   utils.AddMothToTimeForPayment(oldBill.BillEndDate, 1),
	}
	nextBalance := nextBalanceFromBill(oldBill)
	fillDefaultAmountValues(&newBill, nextBalance)
	if error := s.loanRepository.StoreBill(&newBill); error != nil {
		return error
	}
	if error := s.closePeriod(&oldBill); error != nil {
		return error
	}
	return s.recurringBill(loanID)
}

func (s *LoanService) closePeriod(bill *models.Bill) error {
	bill.PeriodStatus = models.PeriodStatusClosed
	return s.loanRepository.UpdateBill(bill)
}

func (s *LoanService) PayLoan(payment *models.Payment) error {
	billsWithDue, error := s.loanRepository.FindBillsWithDueOrOpenOrderedByPeriodAsc(payment.LoanID)
	if error != nil {
		return error
	}
	if error := s.loanRepository.StorePayment(payment); error != nil {
		return error
	}
	remainingPayment := payment.PaymentAmount.RoundBank(config.Round)
	for index, bill := range billsWithDue {
		paymentToBill := decimal.Zero
		if remainingPayment.LessThanOrEqual(decimal.Zero) {
			break
		}
		if remainingPayment.LessThanOrEqual(bill.TotalDue) || len(billsWithDue) == (index+1) {
			paymentToBill = remainingPayment
		} else {
			paymentToBill = bill.TotalDue
		}
		bill.LiquidateBill(payment.PaymentDate)
		if error := s.payBill(paymentToBill, bill); error != nil {
			return error
		}
		remainingPayment = remainingPayment.Sub(paymentToBill).RoundBank(config.Round)
	}
	return nil
}

func (s *LoanService) payBill(paymentToBill decimal.Decimal, bill models.Bill) error {
	billMovement := new(models.BillMovement)
	billMovement.FillInitialBillMovementFromBill(bill)
	bill.ApplyPayment(paymentToBill)
	if bill.FinalPrincipal.LessThanOrEqual(decimal.Zero) {
		bill.PeriodStatus = models.PeriodStatusClosed
		s.closeLoan(bill.LoanID)
	}
	billMovement.FillFinalBillMovementFromBill(bill)
	if error := s.loanRepository.StoreBillMovement(billMovement); error != nil {
		return error
	}
	return s.loanRepository.UpdateBill(&bill)
}

func calculateCloseDateAgreed(loan *models.Loan) {
	loan.CloseDateAgreed = utils.AddMothToTimeForPayment(loan.StartDate, int(loan.PeriodNumbers))
}

func calculatePaymentOfLoan(loan *models.Loan) {
	loan.PaymentAgreed = financial.CalculatePayment(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers)).RoundBank(config.Round)
}

func balanceExpectedInSpecificPeriodOfLoan(loan models.Loan, period int) financial.Balance {
	return financial.BalanceExpectedInSpecificPeriod(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers), period)
}

func nextBalanceFromBill(bill models.Bill) financial.Balance {
	balance := financial.Balance{}
	balance.InitialPrincipal = bill.InitialPrincipal
	balance.Payment = bill.Payment
	balance.InterestRatePeriod = bill.InterestRate
	balance.ToInterest = bill.InterestOfPayment
	balance.ToPrincipal = bill.PrincipalOfPayment
	balance.FinalPrincipal = bill.FinalPrincipal
	return financial.NextBalanceFromBefore(balance)
}

func fillDefaultAmountValues(bill *models.Bill, balance financial.Balance) {
	round := config.Round
	bill.State = models.BillStateDue
	bill.PeriodStatus = models.PeriodStatusOpen
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