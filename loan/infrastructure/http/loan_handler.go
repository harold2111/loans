package http

//TODO: REMOVE CALLS TO MODEL
import (
	"net/http"

	loanApplication "github.com/harold2111/loans/loan/application"

	"github.com/labstack/echo"
)

type HttpLoanHandler struct {
	LoanService loanApplication.LoanService
}

func NewLoanHttpHandler(e *echo.Echo, loanService loanApplication.LoanService) {
	handler := &HttpLoanHandler{
		LoanService: loanService,
	}
	e.GET("/api/loans", handler.handleFindAllLoans)
	e.POST("/api/loans", handler.handleCreateLoan)
	e.POST("/api/loans/simulate", handler.handleSimulateLoan)
	e.POST("/api/loans/pay", handler.handlePayLoan)
}

func (handler *HttpLoanHandler) handleFindAllLoans(context echo.Context) error {
	loanService := handler.LoanService
	response, error := loanService.FindAllLoans()
	if error != nil {
		return error
	}
	return context.JSON(http.StatusOK, response)
}

func (handler *HttpLoanHandler) handleSimulateLoan(context echo.Context) error {
	loanService := handler.LoanService
	var request loanApplication.CreateLoanRequest
	if error := context.Bind(&request); error != nil {
		return error
	}
	response, error := loanService.SimulateLoan(request)
	if error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}

func (handler *HttpLoanHandler) handleCreateLoan(context echo.Context) error {
	loanService := handler.LoanService
	request := loanApplication.CreateLoanRequest{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	if error := loanService.CreateLoan(request); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, request)
}

func (handler *HttpLoanHandler) handlePayLoan(context echo.Context) error {
	loanService := handler.LoanService
	request := loanApplication.PayLoanRequest{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	response, error := loanService.PayLoan(request)
	if error != nil {
		return error
	}
	return context.JSON(http.StatusOK, response)
}
