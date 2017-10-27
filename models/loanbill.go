package models

import (
	"loans/config"
	"loans/errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type LoanBill struct {
	gorm.Model
	LoanID        uint
	State         string
	Period        uint
	BillStartDate time.Time
	BillEndDate   time.Time
	PaymentndDate time.Time
	Payment       decimal.Decimal `gorm:"type:numeric"`
	InterestRate  decimal.Decimal `gorm:"type:numeric"`
	Interest      decimal.Decimal `gorm:"type:numeric"`
	Principal     decimal.Decimal `gorm:"type:numeric"`
	DaysLate      int
	FeeLate       decimal.Decimal `gorm:"type:numeric"`
	Paid          decimal.Decimal `gorm:"type:numeric"`
	ExtraPaid     decimal.Decimal `gorm:"type:numeric"`
	Due           decimal.Decimal `gorm:"type:numeric"`
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

	period := 1
	round := config.Round
	balance := balanceExpectedInSpecificPeriodOfLoan(loan, period)
	newBill := LoanBill{}
	newBill.LoanID = loanID
	newBill.State = "DEUDA"
	newBill.Period = uint(period)
	newBill.BillStartDate = loan.StartDate
	newBill.BillEndDate = addMothToTimeUtil(newBill.BillStartDate, 1)
	newBill.PaymentndDate = newBill.BillEndDate
	newBill.Payment = balance.Payment.RoundBank(round)
	newBill.Interest = balance.ToInterest.RoundBank(round)
	newBill.InterestRate = loan.InterestRatePeriod.RoundBank(round)
	newBill.Principal = balance.ToPrincipal.RoundBank(round)
	newBill.DaysLate = 0
	newBill.FeeLate = decimal.Zero
	newBill.Paid = decimal.Zero
	newBill.ExtraPaid = decimal.Zero
	totalPaid := newBill.Paid.Add(newBill.ExtraPaid)
	totalDue := newBill.Payment.Add(newBill.FeeLate)
	newBill.Due = totalDue.Sub(totalPaid).RoundBank(round)
	newBill.Create()

	return nil

}

func RecurringLoanBillingByLoanID(loanID uint) error {
	loan, error := FindLoanByID(loanID)
	if error != nil {
		return error
	}

	oldLoanBill := LoanBill{}
	config.DB.Raw("SELECT * FROM loan_bills WHERE loan_id = ? AND period = (SELECT max(period) FROM loan_bills where loan_id = ?)", loanID, loanID).Scan(&oldLoanBill)
	now := time.Now().In(oldLoanBill.BillEndDate.Location())
	if now.Before(oldLoanBill.BillEndDate) {
		return nil
	}
	period := int(oldLoanBill.Period + 1)
	round := config.Round
	balance := balanceExpectedInSpecificPeriodOfLoan(loan, period)
	newBill := LoanBill{}
	newBill.LoanID = loanID
	newBill.State = "DEUDA"
	newBill.Period = uint(period)
	newBill.BillStartDate = oldLoanBill.BillEndDate.AddDate(0, 0, 1)
	newBill.BillEndDate = addMothToTimeUtil(oldLoanBill.BillStartDate, 1)
	newBill.PaymentndDate = newBill.BillEndDate
	newBill.Payment = balance.Payment.RoundBank(round)
	newBill.Interest = balance.ToInterest.RoundBank(round)
	newBill.InterestRate = loan.InterestRatePeriod.RoundBank(round)
	newBill.Principal = balance.ToPrincipal.RoundBank(round)
	newBill.DaysLate = 0
	newBill.FeeLate = decimal.Zero
	newBill.Paid = decimal.Zero
	newBill.ExtraPaid = decimal.Zero
	totalPaid := newBill.Paid.Add(newBill.ExtraPaid)
	totalDue := newBill.Payment.Add(newBill.FeeLate)
	newBill.Due = totalDue.Sub(totalPaid).RoundBank(round)
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
