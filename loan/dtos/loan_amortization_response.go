package dtos

type LoanAmortizationsResponse struct {
	LoanResponse
	Amortizations []AmortizationResponse `json:"amortizations"`
}
