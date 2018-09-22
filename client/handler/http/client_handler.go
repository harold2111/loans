package http

import (
	"loans/client"
	"loans/models"
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
	e.DELETE("/api/clients/:id", handler.handleDeleteClient)

	e.GET("/api/clients/:id/addresses", handler.handleFindAddressesByClientID)
	e.POST("/api/clients/:id/addresses", handler.handleCreateAddressClient)
	e.PUT("/api/clients/:id/addresses/:addressID", handler.handleUpdateAddressClient)
	e.DELETE("/api/clients/:id/addresses/:addressID", handler.handleDeletAddressClient)
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
	client.ID = uint(id)
	if error := clientService.UpdateClient(client); error != nil {
		return error
	}
	return context.JSON(http.StatusCreated, client)
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

func (handler *HttpClientHandler) handleFindAddressesByClientID(context echo.Context) error {
	clientService := handler.ClientService
	clientID, _ := strconv.Atoi(context.Param("id"))
	addresses, err := clientService.FindAddressesByClientID(uint(clientID))
	if err != nil {
		return err
	}
	return context.JSON(http.StatusOK, addresses)
}

func (handler *HttpClientHandler) handleCreateAddressClient(context echo.Context) error {
	clientService := handler.ClientService
	clientID, _ := strconv.Atoi(context.Param("id"))
	address := new(models.Address)
	if error := context.Bind(address); error != nil {
		return error
	}
	address.ClientID = uint(clientID)
	err := clientService.CreateAddressClient(address)
	if err != nil {
		return err
	}
	return context.JSON(http.StatusOK, address)
}

func (handler *HttpClientHandler) handleUpdateAddressClient(context echo.Context) error {
	clientService := handler.ClientService
	clientID, _ := strconv.Atoi(context.Param("id"))
	addressID, _ := strconv.Atoi(context.Param("addressID"))
	address := new(models.Address)
	if error := context.Bind(address); error != nil {
		return error
	}
	address.ClientID = uint(clientID)
	address.ID = uint(addressID)
	err := clientService.UpdateAdressClient(address)
	if err != nil {
		return err
	}
	return context.JSON(http.StatusOK, address)
}

func (handler *HttpClientHandler) handleDeletAddressClient(context echo.Context) error {
	clientService := handler.ClientService
	clientID, _ := strconv.Atoi(context.Param("id"))
	addressID, _ := strconv.Atoi(context.Param("addressID"))
	err := clientService.DeleteAddressClient(uint(clientID), uint(addressID))
	if err != nil {
		return err
	}
	return context.JSON(http.StatusAccepted, "OK")
}
