package domain

import (
	"loans/shared/config"
	"time"

	"github.com/shopspring/decimal"
)

type LoanPeriodMovement struct {
	ID                     uint       `gorm:"primary_key" json:"id"`
	CreatedAt              time.Time  `json:"-"`
	UpdatedAt              time.Time  `json:"-"`
	DeletedAt              *time.Time `sql:"index" json:"-"`
	LoandPeriodID          uint
	PaymentID              uint
	MovementDate           time.Time
	InitialPaymentDue      decimal.Decimal `gorm:"type:numeric" json:"initialPaymentDue"`
	InitialFeeLateDue      decimal.Decimal `gorm:"type:numeric" json:"initialFeeLateDue"`
	InitialPaidToPrincipal decimal.Decimal `gorm:"type:numeric" json:"initialPaidToPrincipal"`
	InitialDue             decimal.Decimal `gorm:"type:numeric" json:"initialDue"`
	Paid                   decimal.Decimal `gorm:"type:numeric" json:"paid"`
	DaysLate               int
	PaidToPaymentDue       decimal.Decimal `gorm:"type:numeric" json:"PaidToPaymentDue"`
	PaidToFeeLate          decimal.Decimal `gorm:"type:numeric" json:"PaidToFeeLate"`
	PaidToPrincipal        decimal.Decimal `gorm:"type:numeric" json:"PaidToPrincipal"`
	FinalPaymentDue        decimal.Decimal `gorm:"type:numeric" json:"FinalPaymentDue"`
	FinalFeeLateDue        decimal.Decimal `gorm:"type:numeric" json:"FinalFeeLateDue"`
	FinalDue               decimal.Decimal `gorm:"type:numeric" json:"starFinalDuetDate"`
}

func (billMovement *LoanPeriodMovement) FillInitialBillMovementFromBill(bill LoanPeriod) {
	billMovement.PaymentID = bill.LoanID
	billMovement.BillID = bill.ID
	billMovement.MovementDate = bill.LastLiquidationDate
	billMovement.InitialPaymentDue = bill.PaymentDue
	billMovement.InitialFeeLateDue = bill.FeeLateDue
	billMovement.InitialDue = bill.TotalDue
	billMovement.InitialPaidToPrincipal = bill.PaidToPrincipal
	billMovement.DaysLate = bill.DaysLate
}

func (billMovement *LoanPeriodMovement) FillFinalBillMovementFromBill(bill LoanPeriod) {
	billMovement.PaidToPaymentDue = billMovement.InitialPaymentDue.Sub(bill.PaymentDue).RoundBank(config.Round)
	billMovement.PaidToFeeLate = billMovement.InitialFeeLateDue.Sub(bill.FeeLateDue).RoundBank(config.Round)
	billMovement.FinalPaymentDue = bill.PaymentDue
	billMovement.FinalFeeLateDue = bill.FeeLateDue
	billMovement.FinalDue = bill.TotalDue
	billMovement.PaidToPrincipal = bill.PaidToPrincipal.Sub(billMovement.InitialPaidToPrincipal)
	billMovement.Paid = billMovement.PaidToPaymentDue.Add(billMovement.PaidToFeeLate).Add(billMovement.PaidToPrincipal)
}
