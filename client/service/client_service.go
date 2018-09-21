package service

import (
	"loans/client"
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
	client, err := s.clientRepository.Find(clientID)
	if err != nil {
		return client, err
	}
	return client, nil
}

func (s *clientService) CreateClient(client *models.Client) error {
	if error := s.addDepartmentIDToAddress(client.Addresses); error != nil {
		return error
	}
	return s.clientRepository.Store(client)
}

func (s *clientService) UpdateClient(client *models.Client) error {
	if error := s.addDepartmentIDToAddress(client.Addresses); error != nil {
		return error
	}
	return s.clientRepository.Update(client)
}

func (s *clientService) addDepartmentIDToAddress(addresses []models.Address) error {
	for i, address := range addresses {
		city, err := s.locationRepository.FindCity(address.CityID)
		if err != nil {
			return err
		}
		addresses[i].DepartmentID = city.ID
	}
	return nil
}
