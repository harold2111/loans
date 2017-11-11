package loan

import (
	"loans/dtos"
	"loans/utils"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

func SuscribeLoanHandler(s Service, e *echo.Echo) {
	e.POST("/api/loans", func(c echo.Context) error {
		return handleCreateLoan(s, c)
	})
	e.POST("/api/loans/payments", func(c echo.Context) error {
		return handlePayLoan(s, c)
	})
}

func handleCreateLoan(s Service, context echo.Context) error {
	request := dtos.CreateLoan{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	if error := utils.ValidateStruct(request); error != nil {
		return error
	}
	loan := Loan{}
	if error := copier.Copy(&loan, &request); error != nil {
		return error
	}
	if error := s.CreateLoan(&loan); error != nil {
		return error
	}
	response := dtos.LoanResponse{}
	if error := copier.Copy(&response, &loan); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}

func handlePayLoan(s Service, context echo.Context) error {
	request := dtos.Payment{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	payment := Payment{}
	if error := copier.Copy(&payment, &request); error != nil {
		return error
	}
	if error := s.PayLoan(&payment); error != nil {
		return error
	}
	response := dtos.PaymentResponse{}
	if error := copier.Copy(&response, &payment); error != nil {
		return error
	}
	return context.JSON(http.StatusOK, response)
}
