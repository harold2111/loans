package controllers

import (
	"loans/dtos"
	"loans/models"
	"net/http"

	"github.com/labstack/echo"
)

func PayLoan(context echo.Context) error {
	request := new(dtos.Payment)
	if error := context.Bind(request); error != nil {
		return error
	}
	if error := models.PayLoan(request.LoanID, request.Payment); error != nil {
		return error
	}
	return context.JSON(http.StatusAccepted, request)
}
