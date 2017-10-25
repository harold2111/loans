package models

import (
	"fmt"
	"loans/config"
	"loans/errors"
	"loans/financial"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type LoanBill struct {
	gorm.Model
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

func (loanBill *LoanBill) Create() error {
	error := config.DB.Create(loanBill).Error
	return error
}

func (loanBill *LoanBill) Update() error {
	error := config.DB.Save(loanBill).Error
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
	interestRate := loan.InterestRatePeriod
	newBill.LoanID = loanID
	newBill.State = "ACTIVE"
	newBill.Period = 1
	newBill.BillStartDate = loan.StartDate
	newBill.BillEndDate = addMothToTimeUtil(newBill.BillStartDate, 1)
	newBill.InitialBalance = initialBalance
	newBill.Interest = newBill.InitialBalance.Mul(interestRate).RoundBank(config.Round)
	newBill.DaysLate = 0
	newBill.InterestForLate = decimal.Zero
	newBill.Payment = decimal.Zero
	newBill.ExtraPayment = decimal.Zero
	newBill.PaymentToPrincipal = decimal.Zero
	newBill.FinalBalance = decimal.Zero
	newBill.FinalBalanceExpected =
		financial.BalancingExpectedInSpecificPeriodNumber(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers), int(newBill.Period)).RoundBank(config.Round)
	newBill.Arrrears = decimal.Zero
	newBill.AccumulatedArrears = decimal.Zero
	newBill.Create()

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
	endDatePlusFifteen := oldBill.BillEndDate.AddDate(0, 0, 15)
	//Bill is Fifteen days after las bill payment date
	if oldBill.BillEndDate.After(endDatePlusFifteen) {
		return nil
	}
	var dayLast int
	if !oldBill.Payment.Equal(loan.PaymentAgreed) {
		_, _, dayLast, _, _, _ = diff(oldBill.BillEndDate, endDatePlusFifteen)
	}

	/*************/
	oldBill.State = "CLOSED"
	oldBill.DaysLate = dayLast
	totalPaymentLate := loan.PaymentAgreed.Sub(oldBill.Payment).RoundBank(config.Round)
	annualInterestRatePastDue := financial.EffectiveMonthlyToAnnual(loan.InterestRatePeriod).RoundBank(config.Round)
	fmt.Println("annualInterestRatePastDue: ", annualInterestRatePastDue)
	oldBill.InterestForLate = financial.CalculateInterestPastOfDue(annualInterestRatePastDue, totalPaymentLate, dayLast).RoundBank(config.Round)
	oldBill.FinalBalance = oldBill.InitialBalance.Add(oldBill.Interest).Add(oldBill.InterestForLate).Sub(oldBill.Payment).RoundBank(config.Round)
	oldBill.Arrrears = oldBill.FinalBalance.Sub(oldBill.FinalBalanceExpected).RoundBank(config.Round)
	oldBill.Update()
	/************/

	newBill := &LoanBill{}
	newBill.LoanID = loanID
	newBill.State = "ACTIVE"
	newBill.Period = oldBill.Period + 1
	newBill.BillStartDate = oldBill.BillEndDate
	newBill.BillEndDate = addMothToTimeUtil(loan.StartDate, 1)
	newBill.InitialBalance = oldBill.FinalBalance
	newBill.Interest = newBill.InitialBalance.Mul(loan.InterestRatePeriod).RoundBank(config.Round)
	newBill.DaysLate = 0
	newBill.InterestForLate = decimal.Zero
	newBill.Payment = decimal.Zero
	newBill.ExtraPayment = decimal.Zero
	newBill.FinalBalance = decimal.Zero
	newBill.FinalBalanceExpected =
		financial.BalancingExpectedInSpecificPeriodNumber(
			loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers), int(newBill.Period)).RoundBank(config.Round)
	newBill.Arrrears = decimal.Zero
	newBill.AccumulatedArrears = decimal.Zero
	newBill.Create()
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
