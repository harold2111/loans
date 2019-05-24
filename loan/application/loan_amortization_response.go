package application

type LoanAmortizationsResponse struct {
	LoanResponse
	Amortizations []AmortizationResponse `json:"amortizations"`
}
