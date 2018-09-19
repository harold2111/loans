package client

import "loans/models"

// Repository provides access a client store.
type ClientRepository interface {
	Store(client *models.Client) error
	Update(client *models.Client) error
	Find(clientID uint) (models.Client, error)
	ClientExist(clientID uint) (bool, error)
	FindClientAddress(clientID uint) ([]models.Address, error)
	FindAll() ([]models.Client, error)
}
