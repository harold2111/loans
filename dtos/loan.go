package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

/*
"startDate": "2017-10-30T00:00:00Z",
"closeDate": "2017-10-31T11:14:41-05:00"
*/
type commanLoanFields struct {
	Principal          decimal.Decimal `json:"principal"`
	InterestRatePeriod decimal.Decimal `json:"interestRatePeriod"`
	PeriodNumbers      uint            `json:"periodNumbers"`
	StartDate          time.Time       `json:"startDate"`
	ClientID           uint            `json:"clientID" validate:"required"`
}

type CreateLoan struct {
	commanLoanFields
}

type UpdateLoan struct {
	ID uint `json:"id"`
	commanLoanFields
}

type LoanResponse struct {
	ID uint `json:"id"`
	commanLoanFields
	PaymentAgreed   decimal.Decimal `json:"paymentAgreed"`
	CloseDateAgreed time.Time       `json:"closeDateAgreed"`
	State           string          `json:"state"`
}
