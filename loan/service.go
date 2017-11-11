package loan

import (
	"loans/client"
	"loans/config"
	"loans/errors"
	"loans/utils"
	"time"

	"github.com/shopspring/decimal"
)

type service struct {
	loanRepository   Repository
	clientRepository client.Repository
}

// Service is the interface that provides client methods.
type Service interface {
	CreateLoan(loan *Loan) error
	PayLoan(payment *Payment) error
}

// NewService creates a loan service with necessary dependencies.
func NewService(loanRepository Repository, clientRepository client.Repository) Service {
	return &service{
		loanRepository:   loanRepository,
		clientRepository: clientRepository,
	}
}

func (s *service) CreateLoan(loan *Loan) error {
	loan.State = LoanStateActive
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

func (s *service) closeLoan(loanID uint) error {
	loan, error := s.loanRepository.FindLoanByID(loanID)
	if error != nil {
		return error
	}
	loan.State = LoanStateClosed
	return s.loanRepository.UpdateLoan(&loan)
}

func (s *service) initialBill(loanID uint) error {
	loan, error := s.loanRepository.FindLoanByID(loanID)
	if error != nil {
		return error
	}
	bills, _ := s.loanRepository.FindBillsByLoanID(loanID)
	if len(bills) > 0 {
		return &errors.GracefulError{ErrorCode: errors.BillAlreadyExist}
	}
	period := 1
	newBill := Bill{}
	newBill.LoanID = loan.ID
	newBill.Period = uint(period)
	newBill.BillStartDate = loan.StartDate
	newBill.BillEndDate = utils.AddMothToTimeForPayment(newBill.BillStartDate, 1)
	balancePeriod := balanceExpectedInSpecificPeriodOfLoan(loan, period)
	fillDefaultAmountValues(&newBill, balancePeriod)
	return s.loanRepository.StoreBill(&newBill)
}

func (s *service) recurringBill(loanID uint) error {
	loan, error := s.loanRepository.FindLoanByID(loanID)
	if error != nil {
		return error
	}
	oldBill := Bill{}
	oldBill, error = s.loanRepository.FindBillOpenPeriodByLoanID(loanID)
	if error != nil {
		return error
	}
	if time.Now().Before(oldBill.BillEndDate) {
		return nil
	}
	period := int(oldBill.Period + 1)
	newBill := Bill{}
	newBill.LoanID = loan.ID
	newBill.Period = uint(period)
	newBill.BillStartDate = oldBill.BillEndDate.AddDate(0, 0, 1)
	newBill.BillEndDate = utils.AddMothToTimeForPayment(oldBill.BillEndDate, 1)
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

func (s *service) closePeriod(bill *Bill) error {
	bill.PeriodStatus = PeriodStatusClosed
	return s.loanRepository.UpdateBill(bill)
}

func (s *service) PayLoan(payment *Payment) error {
	billsWithDue, error := s.loanRepository.FindBillsWithDueOrOpenOrderedByPeriodAsc(payment.LoanID)
	if error != nil {
		return error
	}
	if error := s.loanRepository.StorePayment(payment); error != nil {
		return error
	}
	remainingPayment := payment.PaymentAmount.RoundBank(config.Round)
	for index, bill := range billsWithDue {
		if remainingPayment.LessThanOrEqual(decimal.Zero) {
			break
		}
		billMovement := new(BillMovement)
		billMovement.PaymentID = payment.ID
		bill.LiquidateBill(payment.PaymentDate)
		billMovement.fillInitialBillMovementFromBill(bill)

		paymentToBill := decimal.Zero
		if remainingPayment.LessThanOrEqual(bill.TotalDue) || len(billsWithDue) == (index+1) {
			paymentToBill = remainingPayment
		} else {
			paymentToBill = bill.TotalDue
		}
		bill.applyPayment(paymentToBill)
		if bill.FinalPrincipal.LessThanOrEqual(decimal.Zero) {
			bill.PeriodStatus = PeriodStatusClosed
			s.closeLoan(bill.LoanID)
		}
		billMovement.fillFinalBillMovementFromBill(bill)
		if error := s.loanRepository.StoreBillMovement(billMovement); error != nil {
			return error
		}
		if error := s.loanRepository.UpdateBill(&bill); error != nil {
			return error
		}
		remainingPayment = remainingPayment.Sub(paymentToBill).RoundBank(config.Round)
	}
	return nil
}
