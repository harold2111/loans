package domain

import "loans/shared/models"

// Repository provides access a client store.
type ClientRepository interface {
	FindAll() ([]models.Client, error)
	Find(clientID uint) (models.Client, error)
	Create(client *models.Client) error
	Update(client *models.Client) error
	Delete(client *models.Client) error
	FindAddressesByClientID(addressID uint) ([]models.Address, error)
	FindAddressByIDAndClientID(addressID uint, ClientID uint) (*models.Address, error)
	CreateAddressClient(address *models.Address) error
	UpdateAddressClient(address *models.Address) error
	DeleteAddressClient(address *models.Address) error
}