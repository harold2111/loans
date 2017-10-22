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
	ProductID     uint            `json:"productID"`
	Principal     decimal.Decimal `json:"principal"`
	InteresRate   decimal.Decimal `json:"interesRate"`
	PeriodNumbers uint            `json:"periodNumbers"`
	Payment       decimal.Decimal `json:"payment"`
	StartDate     time.Time       `json:"startDate"`
	CloseDate     time.Time       `json:"closeDate"`
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
}
