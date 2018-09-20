package http

import (
	"fmt"
	"loans/client"
	"loans/client/dtos"
	"loans/models"
	"loans/utils"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
)

type HttpClientHandler struct {
	ClientService client.ClientService
}

func NewClientHttpHandler(e *echo.Echo, clientService client.ClientService) {
	handler := &HttpClientHandler{
		ClientService: clientService,
	}
	e.GET("/api/clients", handler.handleFindAllClients)
	e.GET("/api/clients/:id", handler.handleFindClientByID)
	e.POST("/api/clients", handler.handleCreateClient)
	e.PUT("/api/clients/:id", handler.handleUpdateClient)
}

func (handler *HttpClientHandler) handleFindAllClients(context echo.Context) error {
	clientService := handler.ClientService
	clients, error := clientService.FindAllClients()
	if error != nil {
		return error
	}
	return context.JSON(http.StatusOK, clients)
}

func (handler *HttpClientHandler) handleFindClientByID(context echo.Context) error {
	clientService := handler.ClientService
	id, _ := strconv.Atoi(context.Param("id"))
	client, error := clientService.FindClientByID(uint(id))
	if error != nil {
		return error
	}
	response := new(dtos.ClientResponse)
	if error := copier.Copy(response, &client); error != nil {
		return error
	}
	return context.JSON(http.StatusOK, response)
}

func (handler *HttpClientHandler) handleCreateClient(context echo.Context) error {
	clientService := handler.ClientService
	request := new(models.Client)
	if error := context.Bind(request); error != nil {
		return error
	}
	if error := utils.ValidateStruct(request); error != nil {
		return error
	}
	if error := clientService.CreateClient(request); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, request)
}

func (handler *HttpClientHandler) handleUpdateClient(context echo.Context) error {
	clientService := handler.ClientService
	request := new(dtos.UpdateClientRequest)
	id, _ := strconv.Atoi(context.Param("id"))
	if error := context.Bind(request); error != nil {
		return error
	}
	if error := utils.ValidateStruct(request); error != nil {
		return error
	}
	client := new(models.Client)
	if error := copier.Copy(&client, &request); error != nil {
		return error
	}
	client.ID = uint(id)
	fmt.Println(client)
	if error := clientService.UpdateClient(client); error != nil {
		return error
	}
	response := new(dtos.ClientResponse)
	if error := copier.Copy(&response, &client); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, response)
}
