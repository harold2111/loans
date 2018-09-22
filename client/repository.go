package client

import "loans/models"

// Repository provides access a client store.
type ClientRepository interface {
	FindAll() ([]models.Client, error)
	Find(clientID uint) (models.Client, error)
	Create(client *models.Client) error
	Update(client *models.Client) error
	ClientExist(clientID uint) (bool, error)
	FindAddressesByClientID(addressID uint) ([]models.Address, error)
	FindAddressByIDAndClientID(clientID uint, ClientID uint) (models.Address, error)
	CreateClientAddress(address *models.Address) error
	UpdateAddressClient(address *models.Address) error
}
