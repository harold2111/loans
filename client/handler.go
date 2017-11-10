package client

import (
	"loans/dtos"
	"loans/utils"
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
	if error := utils.ValidateStruct(request); error != nil {
		return error
	}
	client := new(Client)
	if error := copier.Copy(&client, &request); error != nil {
		return error
	}
	if error := client.Create(); error != nil {
		return error
	}
	response := new(dtos.ClientResponse)
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
	if error := utils.ValidateStruct(request); error != nil {
		return error
	}
	client := new(Client)
	if error := copier.Copy(&client, &request); error != nil {
		return error
	}
	client.ID = uint(id)
	if error := client.Update(); error != nil {
		return error
	}
	response := new(dtos.ClientResponse)
	if error := copier.Copy(&response, &client); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}
