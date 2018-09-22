package service

import (
	"loans/client"
	"loans/errors"
	"loans/location"
	"loans/models"
	"loans/utils"
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
	if len(client.Addresses) == 0 {
		return &errors.GracefulError{ErrorCode: errors.AtLeastOneAddress}
	}
	for _, address := range client.Addresses {
		if err := s.validateCityID(address.CityID); err != nil {
			return err
		}
	}
	if error := utils.ValidateStruct(client); error != nil {
		return error
	}
	return s.clientRepository.Create(client)
}

func (s *clientService) UpdateClient(client *models.Client) error {
	if exist, err := s.clientExist(client.ID); !exist {
		return err
	}
	if error := utils.ValidateStruct(client); error != nil {
		return error
	}
	return s.clientRepository.Update(client)
}

func (s *clientService) FindAddressesByClientID(clientID uint) ([]models.Address, error) {
	addresses, err := s.clientRepository.FindAddressesByClientID(clientID)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (s *clientService) CreateAddressClient(address *models.Address) error {
	if error := utils.ValidateStruct(address); error != nil {
		return error
	}
	if exist, err := s.clientExist(address.ClientID); !exist {
		return err
	}
	if _, err := s.clientRepository.FindAddressByIDAndClientID(address.ID, address.ClientID); err == nil {
		return &errors.RecordNotFound{ErrorCode: errors.AddressDuplicate}
	}
	if err := s.validateCityID(address.CityID); err != nil {
		return err
	}
	return s.clientRepository.CreateAddressClient(address)
}

func (s *clientService) UpdateAdressClient(address *models.Address) error {
	if error := utils.ValidateStruct(address); error != nil {
		return error
	}
	if _, err := s.clientRepository.FindAddressByIDAndClientID(address.ID, address.ClientID); err != nil {
		return err
	}
	if err := s.validateCityID(address.CityID); err != nil {
		return err
	}
	return s.clientRepository.UpdateAddressClient(address)
}

func (s *clientService) validateCityID(cityID uint) error {
	_, err := s.locationRepository.FindCity(cityID)
	return err
}

func (s *clientService) clientExist(clientID uint) (bool, error) {
	if _, error := s.clientRepository.Find(clientID); error != nil {
		if _, ok := error.(*errors.RecordNotFound); ok {
			return false, error
		}
		return false, error
	}
	return true, nil
}
