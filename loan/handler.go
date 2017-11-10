package loan

import (
	"loans/dtos"
	"loans/utils"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

func CreateLoan(context echo.Context) error {
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
	if error := loan.Create(); error != nil {
		return error
	}
	response := dtos.LoanResponse{}
	if error := copier.Copy(&response, &loan); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}

func PayLoan(context echo.Context) error {
	request := dtos.Payment{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	payment := Payment{}
	if error := copier.Copy(&payment, &request); error != nil {
		return error
	}
	if error := payment.PayLoan(); error != nil {
		return error
	}
	response := dtos.PaymentResponse{}
	if error := copier.Copy(&response, &payment); error != nil {
		return error
	}
	return context.JSON(http.StatusOK, response)
}
