package client

import "loans/models"

// Service is the interface that provides client methods.
type ClientService interface {
	CreateClient(client *models.Client) error
	UpdateClient(client *models.Client) error
	FindAllClients() ([]models.Client, error)
	FindClientByID(clientID uint) (models.Client, error)
}
