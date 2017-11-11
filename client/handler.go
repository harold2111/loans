package client

import (
	"loans/utils"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

func SuscribeClientHandler(s Service, e *echo.Echo) {
	e.POST("/api/clients", func(c echo.Context) error {
		return handleCreateClient(s, c)
	})
	e.PUT("/api/clients/:id", func(c echo.Context) error {
		return handleUpdateClient(s, c)
	})
}

func handleCreateClient(s Service, c echo.Context) error {
	request := new(CreateClientRequest)
	if error := c.Bind(request); error != nil {
		return error
	}
	if error := utils.ValidateStruct(request); error != nil {
		return error
	}
	client := new(Client)
	if error := copier.Copy(&client, &request); error != nil {
		return error
	}
	addresses := new([]Address)
	if error := copier.Copy(addresses, &request.Addresses); error != nil {
		return error
	}
	if error := s.CreateClient(client, addresses); error != nil {
		return error
	}
	response := new(ClientResponse)
	if error := copier.Copy(response, &client); error != nil {
		return error
	}
	return c.JSON(http.StatusCreated, response)
}

func handleUpdateClient(s Service, c echo.Context) error {
	request := new(UpdateClientRequest)
	id, _ := strconv.Atoi(c.Param("id"))
	if error := c.Bind(request); error != nil {
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
	if error := s.UpdateClient(client); error != nil {
		return error
	}
	response := new(ClientResponse)
	if error := copier.Copy(&response, &client); error != nil {
		return error
	}
	return c.JSON(http.StatusCreated, response)
}
