package controllers

import (
	"loans/dtos"
	"loans/models"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

func PayLoan(context echo.Context) error {
	request := dtos.Payment{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	payment := models.Payment{}
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
