package controllers

import (
	"loans/dtos"
	"loans/models"
	"loans/validators"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

func CreateClient(context echo.Context) error {
	request := new(dtos.CreateClient)
	if error := context.Bind(request); error != nil {
		return error
	}
	if error := validators.ValidateStruct(request); error != nil {
		return error
	}
	client := new(models.Client)
	if error := copier.Copy(&client, &request); error != nil {
		return error
	}
	if error := client.Save(); error != nil {
		return error
	}
	response := new(dtos.ClientDTO)
	if error := copier.Copy(&response, &client); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}

func UpdateClient(context echo.Context) error {
	request := new(dtos.UpdateClient)
	id, _ := strconv.Atoi(context.Param("id"))
	if error := context.Bind(request); error != nil {
		return error
	}
	if error := validators.ValidateStruct(request); error != nil {
		return error
	}
	client := new(models.Client)
	if error := copier.Copy(&client, &request); error != nil {
		return error
	}
	client.ID = uint(id)
	if error := client.UpdateClient(); error != nil {
		return error
	}
	response := new(dtos.ClientDTO)
	if error := copier.Copy(&response, &client); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}
