package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	LoanStateActive = "ACTIVE"
	LoanStateClosed = "CLOSED"
)

type Loan struct {
	ID                 uint            `gorm:"primary_key" json:"id"`
	CreatedAt          time.Time       `json:"-"`
	UpdatedAt          time.Time       `json:"-"`
	DeletedAt          *time.Time      `sql:"index" json:"-" `
	Principal          decimal.Decimal `gorm:"type:numeric" json:"principal"`
	InterestRatePeriod decimal.Decimal `gorm:"type:numeric" json:"interestRatePeriod"`
	PeriodNumbers      uint            `json:"periodNumbers"`
	PaymentAgreed      decimal.Decimal `gorm:"type:numeric" json:"paymentAgreed"`
	StartDate          time.Time       `json:"startDate"`
	CloseDateAgreed    time.Time       `json:"closeDateAgreed"`
	CloseDate          *time.Time      `json:"CloseDate"`
	State              string          `json:"state"`
	ClientID           uint            `gorm:"not null" json:"clientID" validate:"required"`
}
