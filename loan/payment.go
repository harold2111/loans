package loan

import (
	"loans/config"
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

func (payment *Payment) PayLoan() error {
	billsWithDue, error := FindBillsWithDueOrOpenOrderedByPeriodAsc(payment.LoanID)
	if error != nil {
		return error
	}
	remainingPayment := payment.PaymentAmount.RoundBank(config.Round)
	payment.Create()
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
		bill.ApplyPayment(paymentToBill)
		billMovement.fillFinalBillMovementFromBill(bill)
		billMovement.Create()
		bill.Update()
		remainingPayment = remainingPayment.Sub(paymentToBill).RoundBank(config.Round)
	}
	return nil

}

func (payment *Payment) Create() error {
	error := config.DB.Create(payment).Error
	return error
}
