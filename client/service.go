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
	CreateClient(client *Client) error
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

func (s *service) CreateClient(client *Client) error {
	if error := validateClientAddress(s.locationRepository, client.Address); error != nil {
		return error
	}	
	return s.clientRepository.Store(client)
}

func (s *service) UpdateClient(client *Client) error {
	if exist, error := s.clientRepository.ClientExist(client.ID); !exist {
		return error
	}
	return s.clientRepository.Update(client)
}

func validateClientAddress(locationRepository location.Repository, address Address) error {
	if len(address.Address) <= 0 {
		return &errors.GracefulError{ErrorCode: errors.AddressRequired}
	}
	if _, error := locationRepository.FindCity(address.CityID); error != nil {
		return error
	}
	return nil
}
