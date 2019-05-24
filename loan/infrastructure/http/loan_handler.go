package http

//TODO: REMOVE CALLS TO MODEL
import (
	loanApplication "loans/loan/application"
	loanDomain "loans/loan/domain"
	"loans/shared/utils"
	"net/http"

	"github.com/jinzhu/copier"
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
	e.POST("/api/loans/simulate", handler.handleSimulateLoan)
	e.POST("/api/loans", handler.handleCreateLoan)
	e.POST("/api/loans/payments", handler.handlePayLoan)
}

func (handler *HttpLoanHandler) handleFindAllLoans(context echo.Context) error {
	loanService := handler.LoanService
	loans, error := loanService.FindAllLoans()
	if error != nil {
		return error
	}
	return context.JSON(http.StatusOK, loans)
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
	if error := utils.ValidateStruct(request); error != nil {
		return error
	}
	loan := loanDomain.Loan{}
	if error := copier.Copy(&loan, &request); error != nil {
		return error
	}
	if error := loanService.CreateLoan(&loan); error != nil {
		return error
	}
	response := loanApplication.LoanResponse{}
	if error := copier.Copy(&response, &loan); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}

func (handler *HttpLoanHandler) handlePayLoan(context echo.Context) error {
	loanService := handler.LoanService
	request := loanApplication.CreatePaymentRequest{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	payment := loanDomain.Payment{}
	if error := copier.Copy(&payment, &request); error != nil {
		return error
	}
	if error := loanService.PayLoan(&payment); error != nil {
		return error
	}
	response := loanApplication.PaymentResponse{}
	if error := copier.Copy(&response, &payment); error != nil {
		return error
	}
	return context.JSON(http.StatusOK, response)
}
