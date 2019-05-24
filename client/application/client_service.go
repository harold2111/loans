package application

import (
	clientDomain "loans/client/domain"
	locationDomain "loans/location/domain"
	"loans/shared/errors"
	"loans/shared/models"
	"loans/shared/utils"
)

type ClientService struct {
	clientRepository   clientDomain.ClientRepository
	locationRepository locationDomain.LocationRepository
}

// NewService creates a client service with necessary dependencies.
func NewClientService(clientRepository clientDomain.ClientRepository, locationRepository locationDomain.LocationRepository) ClientService {
	return ClientService{
		clientRepository:   clientRepository,
		locationRepository: locationRepository,
	}
}

func (s *ClientService) FindAllClients() ([]models.Client, error) {
	return s.clientRepository.FindAll()
}

func (s *ClientService) FindClientByID(clientID uint) (models.Client, error) {
	client, err := s.clientRepository.Find(clientID)
	if err != nil {
		return client, err
	}
	return client, nil
}

func (s *ClientService) CreateClient(client *models.Client) error {
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

func (s *ClientService) UpdateClient(client *models.Client) error {
	if exist, err := s.clientExist(client.ID); !exist {
		return err
	}
	if error := utils.ValidateStruct(client); error != nil {
		return error
	}
	return s.clientRepository.Update(client)
}

func (s *ClientService) DeleteClient(clientID uint) error {
	client, err := s.FindClientByID(clientID)
	if err != nil {
		return err
	}
	return s.clientRepository.Delete(&client)
}

func (s *ClientService) FindAddressesByClientID(clientID uint) ([]models.Address, error) {
	addresses, err := s.clientRepository.FindAddressesByClientID(clientID)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (s *ClientService) CreateAddressClient(address *models.Address) error {
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

func (s *ClientService) UpdateAdressClient(address *models.Address) error {
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

func (s *ClientService) DeleteAddressClient(clientID uint, addressID uint) error {
	address, err := s.clientRepository.FindAddressByIDAndClientID(addressID, clientID)
	if err != nil {
		return err
	}
	return s.clientRepository.DeleteAddressClient(address)
}

func (s *ClientService) FindAddressByClientIDAndAddressID(addressID uint, clientID uint) (*models.Address, error) {
	return s.clientRepository.FindAddressByIDAndClientID(addressID, clientID)
}

func (s *ClientService) validateCityID(cityID uint) error {
	_, err := s.locationRepository.FindCity(cityID)
	return err
}

func (s *ClientService) clientExist(clientID uint) (bool, error) {
	if _, error := s.clientRepository.Find(clientID); error != nil {
		if _, ok := error.(*errors.RecordNotFound); ok {
			return false, error
		}
		return false, error
	}
	return true, nil
}
