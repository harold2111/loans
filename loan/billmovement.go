package loan

import (
	"loans/config"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type BillMovement struct {
	gorm.Model
	BillID                 uint
	PaymentID              uint
	MovementDate           time.Time
	InitialPaymentDue      decimal.Decimal `gorm:"type:numeric"`
	InitialFeeLateDue      decimal.Decimal `gorm:"type:numeric"`
	InitialPaidToPrincipal decimal.Decimal `gorm:"type:numeric"`
	InitialDue             decimal.Decimal `gorm:"type:numeric"`
	Paid                   decimal.Decimal `gorm:"type:numeric"`
	DaysLate               int
	PaidToPaymentDue       decimal.Decimal `gorm:"type:numeric"`
	PaidToFeeLate          decimal.Decimal `gorm:"type:numeric"`
	PaidToPrincipal        decimal.Decimal `gorm:"type:numeric"`
	FinalPaymentDue        decimal.Decimal `gorm:"type:numeric"`
	FinalFeeLateDue        decimal.Decimal `gorm:"type:numeric"`
	FinalDue               decimal.Decimal `gorm:"type:numeric"`
}

func (billMovement *BillMovement) fillInitialBillMovementFromBill(bill Bill) {
	billMovement.BillID = bill.ID
	billMovement.MovementDate = bill.LastLiquidationDate
	billMovement.InitialPaymentDue = bill.PaymentDue
	billMovement.InitialFeeLateDue = bill.FeeLateDue
	billMovement.InitialDue = bill.TotalDue
	billMovement.InitialPaidToPrincipal = bill.PaidToPrincipal
	billMovement.DaysLate = bill.DaysLate
}

func (billMovement *BillMovement) fillFinalBillMovementFromBill(bill Bill) {
	billMovement.PaidToPaymentDue = billMovement.InitialPaymentDue.Sub(bill.PaymentDue).RoundBank(config.Round)
	billMovement.PaidToFeeLate = billMovement.InitialFeeLateDue.Sub(bill.FeeLateDue).RoundBank(config.Round)
	billMovement.FinalPaymentDue = bill.PaymentDue
	billMovement.FinalFeeLateDue = bill.FeeLateDue
	billMovement.FinalDue = bill.TotalDue
	billMovement.PaidToPrincipal = bill.PaidToPrincipal.Sub(billMovement.InitialPaidToPrincipal)
	billMovement.Paid = billMovement.PaidToPaymentDue.Add(billMovement.PaidToFeeLate).Add(billMovement.PaidToPrincipal)
}
