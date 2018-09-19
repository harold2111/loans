package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

const (
	LoanStateActive = "ACTIVE"
	LoanStateClosed = "CLOSED"
)

type Loan struct {
	gorm.Model
	Principal          decimal.Decimal `gorm:"type:numeric"`
	InterestRatePeriod decimal.Decimal `gorm:"type:numeric"`
	PeriodNumbers      uint
	PaymentAgreed      decimal.Decimal `gorm:"type:numeric"`
	StartDate          time.Time
	CloseDateAgreed    time.Time
	CloseDate          *time.Time
	State              string
	ClientID           uint `gorm:"not null"`
}
