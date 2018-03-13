package client

import (
	"loans/errors"
	"loans/location"
)

type service struct {
	clientRepository   Repository
	locationRepository location.Repository
}

// Service is the interface that provides client methods.
type Service interface {
	CreateClient(client *Client, adresses *[]Address) error
	UpdateClient(client *Client) error
	FindAllClients() ([]Client, error)
}

// NewService creates a client service with necessary dependencies.
func NewService(clientRepository Repository, locationRepository location.Repository) Service {
	return &service{
		clientRepository:   clientRepository,
		locationRepository: locationRepository,
	}
}

func (s *service) FindAllClients() ([]Client, error) {
	return s.clientRepository.FindAll()
}

func (s *service) CreateClient(client *Client, adresses *[]Address) error {
	if error := validateClientAddress(s.locationRepository, adresses); error != nil {
		return error
	}
	if error := s.clientRepository.Store(client); error != nil {
		return error
	}
	return s.clientRepository.StoreClientAddresses(client.ID, adresses)
}

func (s *service) UpdateClient(client *Client) error {
	if exist, error := s.clientRepository.ClientExist(client.ID); !exist {
		return error
	}
	return s.clientRepository.Update(client)
}

func validateClientAddress(locationRepository location.Repository, addresses *[]Address) error {
	if len(*addresses) <= 0 {
		return &errors.GracefulError{ErrorCode: errors.AddressRequired}
	}
	for _, addressValue := range *addresses {
		if _, error := locationRepository.FindCity(addressValue.CityID); error != nil {
			return error
		}
	}
	return nil
}
