package domain

import (
	"time"

	"github.com/harold2111/loans/shared/errors"
	"github.com/harold2111/loans/shared/utils"
	"github.com/shopspring/decimal"
)

const (

	//ExtraToPrincipal if the payment has extra amount should be applied to the principal.
	ExtraToPrincipal = "ExtraToPrincipal"
	//ExtraToNextPeriods if the payment has extra amount should be applied to the next periods.
	ExtraToNextPeriods = "ExtraToNextPeriods"
)

//Payment represents a loan payment
type Payment struct {
	ID              int `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time      `sql:"index"`
	LoanID          int             `gorm:"not null" validate:"required"`
	PaymentAmount   decimal.Decimal `gorm:"not null; type:numeric"`
	RemainingAmount decimal.Decimal `gorm:"not null; type:numeric"`
	PaymentDate     time.Time
	PaymentType     string
}

//NewPayment create a new payment
func NewPayment(loanID int, paymentAmount decimal.Decimal, paymentDate time.Time, paymentType string) (Payment, error) {
	payment := Payment{
		LoanID:          loanID,
		PaymentAmount:   paymentAmount,
		RemainingAmount: paymentAmount,
		PaymentDate:     paymentDate,
		PaymentType:     paymentType,
	}
	error := payment.validateForCreation()
	if error != nil {
		return Payment{}, error
	}
	return payment, nil

}

func (p *Payment) validateForCreation() error {
	if error := utils.ValidateVar("loanID", p.LoanID, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("paymentAmount", p.PaymentAmount, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("paymentDate", p.PaymentDate, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("PaymentType", p.PaymentType, "required"); error != nil {
		return error
	} else if !p.isExtraToNextPeriods() && !p.isExtraToPrincipal() {
		return &errors.GracefulError{ErrorCode: errors.PaymentTypeInvalid}
	}
	return nil
}

func (p Payment) isExtraToPrincipal() bool {
	return p.PaymentType == ExtraToPrincipal
}

func (p Payment) isExtraToNextPeriods() bool {
	return p.PaymentType == ExtraToNextPeriods
}
