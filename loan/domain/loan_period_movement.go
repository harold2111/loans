package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type LoanPeriodMovement struct {
	ID              uint `gorm:"primary_key"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `sql:"index" `
	LoandPeriodID   uint
	PaymentID       uint
	LiquidationDate time.Time

	InitialDebtForArrears decimal.Decimal `gorm:"type:numeric"`
	InitialDebtOfPayment  decimal.Decimal `gorm:"type:numeric"`
	InitialDebt           decimal.Decimal `gorm:"type:numeric"`

	Paid                               decimal.Decimal `gorm:"type:numeric"`
	DaysInArrearsSinceLastLiquidation  int
	DebtForArrearsSinceLastLiquidation decimal.Decimal `gorm:"type:numeric"`
	PaidToDebtForArrears               decimal.Decimal `gorm:"type:numeric"`
	PaidToPaymentDebt                  decimal.Decimal `gorm:"type:numeric"`
	PaidExtraToPrincipal               decimal.Decimal `gorm:"type:numeric"`

	FinalDebtForArrears decimal.Decimal `gorm:"type:numeric"`
	FinalDebtOfPayment  decimal.Decimal `gorm:"type:numeric"`
	FinalDebt           decimal.Decimal `gorm:"type:numeric"`
}

func (loanPeriodMovement *LoanPeriodMovement) fillInitialMovementFromPeriod(period LoanPeriod) {
	loanPeriodMovement.LoandPeriodID = period.ID
	loanPeriodMovement.LiquidationDate = period.LastLiquidationDate
	loanPeriodMovement.DaysInArrearsSinceLastLiquidation = period.DaysInArrearsSinceLastLiquidation
	loanPeriodMovement.DebtForArrearsSinceLastLiquidation = period.DebtForArrearsSinceLastLiquidation
	loanPeriodMovement.InitialDebtForArrears = period.TotalDebtForArrears
	loanPeriodMovement.InitialDebtOfPayment = period.TotalDebtOfPayment
	loanPeriodMovement.InitialDebt = period.TotalDebt
}

func (loanPeriodMovement *LoanPeriodMovement) fillFinalMovementFromPeriod(period LoanPeriod) {
	loanPeriodMovement.FinalDebtForArrears = period.TotalDebtForArrears
	loanPeriodMovement.FinalDebtOfPayment = period.TotalDebtOfPayment
	loanPeriodMovement.FinalDebt = period.TotalDebt
}