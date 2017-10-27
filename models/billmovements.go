package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type BillMovement struct {
	gorm.Model
	MovementDate   time.Time
	MovementType   string
	Payment        decimal.Decimal
	InitialBalance decimal.Decimal
	FinalBalance   decimal.Decimal
	Note           string
}
