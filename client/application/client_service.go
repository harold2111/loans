package application

import (
	clientDomain "github.com/harold2111/loans/client/domain"
	locationDomain "github.com/harold2111/loans/location/domain"
	"github.com/harold2111/loans/shared/errors"
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

func (s *ClientService) CreateClient(createClientRequest CreateClientRequest) error {
	var clientAddresses []clientDomain.Address
	for _, createClientAddress := range createClientRequest.Addresses {
		address, err := clientDomain.NewAddessForCreateClient(
			createClientAddress.StretAddress,
			createClientAddress.DepartmentID,
			createClientAddress.CityID)
		if err != nil {
			return err
		}
		clientAddresses = append(clientAddresses, address)
	}
	client, error := clientDomain.NewClientForCreate(
		createClientRequest.Identification,
		createClientRequest.FirstName,
		createClientRequest.LastName,
		createClientRequest.Telephone1,
		createClientRequest.Telephone2,
		createClientRequest.Email,
		clientAddresses,
	)
	if error != nil {
		return error
	}
	for _, address := range createClientRequest.Addresses {
		if err := s.validateCityID(address.CityID); err != nil {
			return err
		}
	}
	return s.clientRepository.Create(&client)
}

func (s *ClientService) FindAllClients() ([]GetClientResponse, error) {
	clients, err := s.clientRepository.FindAll()
	var getClientsResponse []GetClientResponse
	if err != nil {
		return getClientsResponse, err
	}
	for _, client := range clients {
		var getClientResponse GetClientResponse
		getClientResponse.fillFromClient(client)
		getClientsResponse = append(getClientsResponse, getClientResponse)
	}
	return getClientsResponse, nil
}

func (s *ClientService) FindClientByID(clientID uint) (GetClientResponse, error) {
	client, err := s.clientRepository.Find(clientID)
	var getClientResponse GetClientResponse
	if err != nil {
		return getClientResponse, err
	}
	getClientResponse.fillFromClient(client)
	return getClientResponse, nil
}

func (s *ClientService) UpdateClient(updateClientRequest UpdateClientRequest) error {
	var clientAddresses []clientDomain.Address
	clientID := updateClientRequest.ID
	if exist, err := s.clientExist(clientID); !exist {
		return err
	}
	for _, updateAddressClientRequest := range updateClientRequest.Addresses {
		clientAddress, err := clientDomain.NewAddessForUpdateClient(
			updateAddressClientRequest.ID,
			updateAddressClientRequest.StretAddress,
			updateAddressClientRequest.DepartmentID,
			updateAddressClientRequest.CityID,
		)
		if err != nil {
			return err
		}
		clientAddresses = append(clientAddresses, clientAddress)
	}
	client, err := clientDomain.NewClientForUpdate(
		clientID,
		updateClientRequest.Identification,
		updateClientRequest.FirstName,
		updateClientRequest.LastName,
		updateClientRequest.Telephone1,
		updateClientRequest.Telephone2,
		updateClientRequest.Email,
		clientAddresses,
	)
	if err != nil {
		return err
	}
	return s.clientRepository.Update(&client)
}

func (s *ClientService) DeleteClient(clientID uint) error {
	client, err := s.clientRepository.Find(clientID)
	if err != nil {
		return err
	}
	return s.clientRepository.Delete(&client)
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
