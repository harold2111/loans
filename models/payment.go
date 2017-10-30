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
