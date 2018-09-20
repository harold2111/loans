package service

import (
	"loans/client"
	"loans/errors"
	"loans/location"
	"loans/models"
)

type clientService struct {
	clientRepository   client.ClientRepository
	locationRepository location.LocationRepository
}

// NewService creates a client service with necessary dependencies.
func NewClientService(clientRepository client.ClientRepository, locationRepository location.LocationRepository) client.ClientService {
	return &clientService{
		clientRepository:   clientRepository,
		locationRepository: locationRepository,
	}
}

func (s *clientService) FindAllClients() ([]models.Client, error) {
	return s.clientRepository.FindAll()
}

func (s *clientService) FindClientByID(clientID uint) (models.Client, error) {
	var client models.Client
	var addresses []models.Address
	var err error

	client, err = s.clientRepository.Find(clientID)
	if err != nil {
		return client, err
	}
	addresses, err = s.clientRepository.FindClientAddress(clientID)
	if err != nil {
		return client, err
	}
	client.Addresses = addresses
	return client, nil
}

func (s *clientService) CreateClient(client *models.Client) error {
	//TODO: valide all address
	if error := validateClientAddress(s.locationRepository, client.Addresses[0]); error != nil {
		return error
	}
	return s.clientRepository.Store(client)
}

func (s *clientService) UpdateClient(client *models.Client) error {
	if exist, error := s.clientRepository.ClientExist(client.ID); !exist {
		return error
	}
	return s.clientRepository.Update(client)
}

func validateClientAddress(locationRepository location.LocationRepository, address models.Address) error {
	if len(address.Address) <= 0 {
		return &errors.GracefulError{ErrorCode: errors.AddressRequired}
	}
	if _, error := locationRepository.FindCity(address.CityID); error != nil {
		return error
	}
	return nil
}
