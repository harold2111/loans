package http

import (
	clientApplication "loans/client/application"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type HttpClientHandler struct {
	ClientService clientApplication.ClientService
}

func NewClientHttpHandler(e *echo.Echo, clientService clientApplication.ClientService) {
	handler := &HttpClientHandler{
		ClientService: clientService,
	}
	e.POST("/api/clients", handler.handleCreateClient)
	e.GET("/api/clients", handler.handleFindAllClients)
	e.GET("/api/clients/:id", handler.handleFindClientByID)
	e.PUT("/api/clients/:id", handler.handleUpdateClient)
	e.DELETE("/api/clients/:id", handler.handleDeleteClient)
}

func (handler *HttpClientHandler) handleCreateClient(context echo.Context) error {
	clientService := handler.ClientService
	request := clientApplication.CreateClientRequest{}
	if error := context.Bind(&request); error != nil {
		return error
	}
	if error := clientService.CreateClient(request); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, request)
}

func (handler *HttpClientHandler) handleFindAllClients(context echo.Context) error {
	clientService := handler.ClientService
	response, error := clientService.FindAllClients()
	if error != nil {
		return error
	}
	return context.JSON(http.StatusOK, response)
}

func (handler *HttpClientHandler) handleFindClientByID(context echo.Context) error {
	clientService := handler.ClientService
	id, _ := strconv.Atoi(context.Param("id"))
	getClientResponse, error := clientService.FindClientByID(uint(id))
	if error != nil {
		return error
	}
	return context.JSON(http.StatusOK, getClientResponse)
}

func (handler *HttpClientHandler) handleUpdateClient(context echo.Context) error {
	request := clientApplication.UpdateClientRequest{}
	id, _ := strconv.Atoi(context.Param("id"))
	if error := context.Bind(&request); error != nil {
		return error
	}
	request.ID = uint(id)
	clientService := handler.ClientService
	if error := clientService.UpdateClient(request); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, request)
}

func (handler *HttpClientHandler) handleDeleteClient(context echo.Context) error {
	clientService := handler.ClientService
	clientID, _ := strconv.Atoi(context.Param("id"))
	err := clientService.DeleteClient(uint(clientID))
	if err != nil {
		return err
	}
	return context.JSON(http.StatusAccepted, "OK")
}
