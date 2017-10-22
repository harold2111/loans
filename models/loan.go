package models

import (
	"loans/config"
	"time"

	"github.com/shopspring/decimal"

	"github.com/jinzhu/gorm"
)

const (
	LOAN = "LOAN"
)

type Product struct {
	gorm.Model
	ProductType string
}

type Loan struct {
	gorm.Model
	ProductID     uint
	Principal     decimal.Decimal `gorm:"type:numeric"`
	InteresRate   decimal.Decimal `gorm:"type:numeric"`
	PeriodNumbers uint
	Payment       decimal.Decimal `gorm:"type:numeric"`
	StartDate     time.Time
	CloseDate     time.Time
}

func (loan *Loan) Save() error {
	error := config.DB.Create(loan).Error
	return error
}
