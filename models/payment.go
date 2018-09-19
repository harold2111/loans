package models

import (
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
