package controllers

import (
	"loans/dtos"
	"loans/models"
	"loans/validators"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

func CreateLoan(context echo.Context) error {
	request := new(dtos.CreateLoan)
	if error := context.Bind(request); error != nil {
		return error
	}
	if error := validators.ValidateStruct(request); error != nil {
		return error
	}
	loan := new(models.Loan)
	if error := copier.Copy(&loan, &request); error != nil {
		return error
	}
	if error := loan.Create(); error != nil {
		return error
	}
	response := new(dtos.LoanResponse)
	if error := copier.Copy(&response, &loan); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}
