package client

import "loans/models"

// Service is the interface that provides client methods.
type ClientService interface {
	FindAllClients() ([]models.Client, error)
	FindClientByID(clientID uint) (models.Client, error)
	CreateClient(client *models.Client) error
	UpdateClient(client *models.Client) error
	DeleteClient(clientID uint) error
	FindAddressesByClientID(clientID uint) ([]models.Address, error)
	CreateAddressClient(address *models.Address) error
	UpdateAdressClient(address *models.Address) error
	DeleteAddressClient(clientID uint, addressID uint) error
	FindAddressByClientIDAndAddressID(addressID uint, clientID uint) (*models.Address, error)
}
