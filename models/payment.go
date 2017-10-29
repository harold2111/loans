package models

import (
	"loans/config"

	"github.com/shopspring/decimal"
)

func PayLoan(loanID uint, payment decimal.Decimal) error {
	billsWithDue, error := FindBillsWithDueOrOpenOrderedByPeriodAsc(loanID)
	if error != nil {
		return error
	}
	remainingPayment := payment.RoundBank(config.Round)
	for index, bill := range billsWithDue {
		if remainingPayment.LessThanOrEqual(decimal.Zero) {
			break
		}
		billMovement := new(BillMovement)
		bill.LiquidateBill()
		billMovement.fillInitalBillMovementFromBill(bill)

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

func (bill *Bill) ApplyPayment(paymentToBill decimal.Decimal) {
	//the payment NO covers all the fee late
	if paymentToBill.LessThanOrEqual(bill.FeeLateDue) {
		bill.FeeLateDue = bill.FeeLateDue.Sub(paymentToBill)
	} else { //the payment covers fee late
		remainingPaymentToBill := paymentToBill.Sub(bill.FeeLateDue)
		bill.FeeLateDue = decimal.Zero
		paymentDue := bill.PaymentDue.Sub(remainingPaymentToBill).RoundBank(config.Round)
		if paymentDue.LessThanOrEqual(decimal.Zero) {
			bill.PaidToPrincipal = bill.PaidToPrincipal.Add(paymentDue.Abs()).RoundBank(config.Round)
			bill.PaymentDue = decimal.Zero
		} else {
			bill.PaymentDue = paymentDue
		}
	}
	bill.TotalDue = bill.PaymentDue.Add(bill.FeeLateDue).RoundBank(config.Round)
	bill.Paid = bill.Paid.Add(paymentToBill)
	if bill.TotalDue.LessThanOrEqual(decimal.Zero) {
		bill.State = BillStatePaid
	}
}
