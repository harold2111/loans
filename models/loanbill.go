package models

import (
	"loans/config"
	"loans/errors"
	"time"

	"github.com/shopspring/decimal"
)

type LoanBill struct {
	LoanID               uint
	State                string
	Period               uint
	BillStartDate        time.Time
	BillEndDate          time.Time
	InitialBalance       decimal.Decimal `gorm:"type:numeric"`
	Interest             decimal.Decimal `gorm:"type:numeric"`
	DaysLate             int
	InterestForLate      decimal.Decimal `gorm:"type:numeric"`
	Payment              decimal.Decimal `gorm:"type:numeric"`
	ExtraPayment         decimal.Decimal `gorm:"type:numeric"`
	PaymentToPrincipal   decimal.Decimal `gorm:"type:numeric"`
	FinalBalance         decimal.Decimal `gorm:"type:numeric"`
	FinalBalanceExpected decimal.Decimal `gorm:"type:numeric"`
	Arrrears             decimal.Decimal `gorm:"type:numeric"`
	AccumulatedArrears   decimal.Decimal `gorm:"type:numeric"`
}

func (loanBill *LoanBill) Save() error {
	error := config.DB.Create(loanBill).Error
	return error
}

func CreateInitialBill(loanID uint) error {
	loan, error := FindLoanByID(loanID)
	if error != nil {
		return error
	}
	var loanBills []LoanBill
	config.DB.Find(&loanBills, "loan_id = ?", loanID)
	if len(loanBills) > 0 {
		return &errors.GracefulError{ErrorCode: errors.BillAlreadyExist}
	}

	var newBill LoanBill
	initialBalance := loan.Principal
	interestRate := loan.InterestRate
	newBill.LoanID = loanID
	newBill.State = "ACTIVE"
	newBill.Period = 1
	newBill.BillStartDate = loan.StartDate
	newBill.BillEndDate = addMothToTimeUtil(newBill.BillStartDate, 1)
	newBill.InitialBalance = initialBalance
	newBill.Interest = newBill.InitialBalance.Mul(interestRate.Div(decimal.NewFromFloat(100))).RoundBank(5)
	newBill.DaysLate = 0
	newBill.InterestForLate = decimal.Zero
	newBill.Payment = decimal.Zero
	newBill.ExtraPayment = decimal.Zero
	newBill.PaymentToPrincipal = decimal.Zero
	newBill.FinalBalance = decimal.Zero
	newBill.FinalBalanceExpected =
		getBalancingInSpecificPeriodNumber(loan.Principal, loan.InterestRate, int(loan.PeriodNumbers), int(newBill.Period))
	newBill.Arrrears = decimal.Zero
	newBill.AccumulatedArrears = decimal.Zero
	newBill.Save()

	return nil

}

func RecurringLoanBillingByLoanID(loanID uint) error {
	loan, error := FindLoanByID(loanID)
	if error != nil {
		return error
	}
	var loanBills []LoanBill
	config.DB.Find(&loanBills, "loan_id = ? AND state = ?", loanID, "ACTIVE")
	if len(loanBills) != 1 {
		return &errors.GracefulError{ErrorCode: errors.ToManyBillActives}
	}
	oldBill := loanBills[0]
	now := time.Now().In(oldBill.BillEndDate.Location())
	if oldBill.BillEndDate.After(now) {
		return nil
	}
	var dayLast int
	if !oldBill.Payment.Equal(loan.PaymentAgreed) {
		_, _, dayLast, _, _, _ = diff(oldBill.BillEndDate, now)
	}

	/*************/
	oldBill.State = "CLOSED"
	oldBill.DaysLate = dayLast
	totalPaymentLate := loan.PaymentAgreed.Sub(oldBill.Payment).RoundBank(5)
	oldBill.InterestForLate = CalculateInterestPastOfDue(loan.InterestRate.Div(decimal.NewFromFloat(100)), totalPaymentLate, dayLast)
	oldBill.FinalBalance = oldBill.InitialBalance.Add(oldBill.Interest).Add(oldBill.InterestForLate).Sub(oldBill.Payment).RoundBank(5)
	oldBill.Save()
	/************/

	var newBill LoanBill
	initialBalance := oldBill.FinalBalance
	interestRate := loan.InterestRate
	newBill.LoanID = loanID
	newBill.State = "ACTIVE"
	newBill.Period = oldBill.Period + 1
	newBill.BillStartDate = oldBill.BillEndDate
	newBill.BillEndDate = addMothToTimeUtil(loan.StartDate, 1)
	newBill.InitialBalance = initialBalance
	newBill.Interest = newBill.InitialBalance.Mul(interestRate.Div(decimal.NewFromFloat(100))).RoundBank(5)
	newBill.DaysLate = dayLast
	newBill.InterestForLate = decimal.Zero
	newBill.Payment = decimal.Zero
	newBill.ExtraPayment = decimal.Zero
	newBill.FinalBalance = decimal.Zero
	newBill.FinalBalanceExpected =
		getBalancingInSpecificPeriodNumber(loan.Principal, loan.InterestRate, int(loan.PeriodNumbers), int(newBill.Period))
	newBill.Arrrears = newBill.FinalBalance.Sub(newBill.FinalBalanceExpected).RoundBank(5)
	newBill.AccumulatedArrears = decimal.Zero
	newBill.Save()
	return nil
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
