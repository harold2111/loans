package loan

import (
	"loans/loan"
	"loans/loan/dtos"
	"loans/models"
	"loans/utils"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

type HttpLoanHandler struct {
	LoanService loan.LoanService
}

func NewLoanHttpHandler(e *echo.Echo, loanService loan.LoanService) {
	handler := &HttpLoanHandler{
		LoanService: loanService,
	}
	e.POST("/api/loans", handler.handleCreateLoan)
	e.POST("/api/loans/payments", handler.handlePayLoan)
}

func (handler *HttpLoanHandler) handleCreateLoan(context echo.Context) error {
	loanService := handler.LoanService
	request := dtos.CreateLoan{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	if error := utils.ValidateStruct(request); error != nil {
		return error
	}
	loan := models.Loan{}
	if error := copier.Copy(&loan, &request); error != nil {
		return error
	}
	if error := loanService.CreateLoan(&loan); error != nil {
		return error
	}
	response := dtos.LoanResponse{}
	if error := copier.Copy(&response, &loan); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}

func (handler *HttpLoanHandler) handlePayLoan(context echo.Context) error {
	loanService := handler.LoanService
	request := dtos.PaymentRequest{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	payment := models.Payment{}
	if error := copier.Copy(&payment, &request); error != nil {
		return error
	}
	if error := loanService.PayLoan(&payment); error != nil {
		return error
	}
	response := dtos.PaymentResponse{}
	if error := copier.Copy(&response, &payment); error != nil {
		return error
	}
	return context.JSON(http.StatusOK, response)
}