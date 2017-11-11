package loan

import (
	"loans/config"
	"loans/financial"
	"loans/utils"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/shopspring/decimal"
)

type Payment struct {
	gorm.Model
	LoanID        uint
	PaymentAmount decimal.Decimal `gorm:"type:numeric"`
	PaymentDate   time.Time
}

func calculateCloseDateAgreed(loan *Loan) {
	loan.CloseDateAgreed = utils.AddMothToTimeForPayment(loan.StartDate, int(loan.PeriodNumbers))
}

func calculatePaymentOfLoan(loan *Loan) {
	loan.PaymentAgreed = financial.CalculatePayment(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers)).RoundBank(config.Round)

}

func balanceExpectedInSpecificPeriodOfLoan(loan Loan, period int) financial.Balance {
	return financial.BalanceExpectedInSpecificPeriod(loan.Principal, loan.InterestRatePeriod, int(loan.PeriodNumbers), period)
}
