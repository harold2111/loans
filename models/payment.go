package models

import (
	"loans/config"
	"loans/errors"
	"time"

	"github.com/shopspring/decimal"
)

type LoanBill struct {
	LoanID                     uint
	State                      string
	Period                     uint
	BillStartDate              time.Time
	BillEndDate                time.Time
	InitialBalance             decimal.Decimal
	Interest                   decimal.Decimal
	DaysLate                   int
	InterestForLate            decimal.Decimal
	Payment                    decimal.Decimal
	ExtraPayment               decimal.Decimal
	PaymentToPrincipal         decimal.Decimal
	FinalBalance               decimal.Decimal
	FinalBalanceWithoutArrears decimal.Decimal
	Arrrears                   decimal.Decimal
	AccumulatedArrears         decimal.Decimal
}

func (loanBill *LoanBill) CreateLoanBill() error {
	loan, error := FindLoanByID(loanBill.LoanID)
	if error != nil {
		return error
	}
	var loans []Loan
	config.DB.Find(&loans, "loan_id = ? AND state IN (?)", loanBill.LoanID, []string{"ACTIVE", "CLOSED"})
	if len(loans) > 1 {
		return &errors.GracefulError{ErrorCode: errors.ToManyBillActives}
	}

	if len(loans) > 0 {
		initialBalance := loan.Principal
		interestRate := loan.InterestRate
		paymentAgreed := loan.PaymentAgreed
		loanBill.State = "ACTIVE"
		loanBill.Period = 1
		loanBill.BillStartDate = loan.StartDate
		loanBill.BillEndDate = addMothToTimeUtil(loan.StartDate, 1)
		loanBill.InitialBalance = initialBalance
		loanBill.Interest = loanBill.InitialBalance.Mul(interestRate)
		loanBill.DaysLate = 0
		loanBill.InterestForLate = decimal.Zero
		loanBill.Payment = decimal.Zero
		loanBill.ExtraPayment = decimal.Zero
		loanBill.FinalBalance = loanBill.InitialBalance.Add(loanBill.Interest).Add(loanBill.InterestForLate).Sub(loanBill.Payment)
		loanBill.FinalBalanceWithoutArrears =
			getBalancingInSpecificPeriodNumber(loan.Principal, loan.InterestRate, int(loan.PeriodNumbers), int(loanBill.Period))

	} else {
		//TODO: crear el que sigue
	}
	return nil
}
