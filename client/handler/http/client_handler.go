package http

import (
	"loans/client"
	"loans/models"
	"loans/utils"
	"net/http"
	"strconv"

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
	return context.JSON(http.StatusOK, client)
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
	client := new(models.Client)
	id, _ := strconv.Atoi(context.Param("id"))
	if error := context.Bind(client); error != nil {
		return error
	}
	if error := utils.ValidateStruct(client); error != nil {
		return error
	}
	client.ID = uint(id)
	if error := clientService.UpdateClient(client); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, client)
}
