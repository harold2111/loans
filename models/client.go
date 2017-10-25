package models

import (
	"fmt"
	"loans/config"
	"loans/errors"
	"loans/postgres"

	"github.com/jinzhu/gorm"
)

type Client struct {
	gorm.Model
	Identification string `gorm:"not null; unique_index"`
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	Telephone1     string `gorm:"not null"`
	Telephone2     string
	Addresses      []Address
}

const (
	UniqueConstraintIdentification = "uix_clients_identification"
)

func (client *Client) Create() error {
	if error := validateClientAddress(client.Addresses); error != nil {
		fmt.Println(error)
		return error
	}
	error := config.DB.Create(client).Error
	if error != nil {
		if postgres.IsUniqueConstraintError(error, UniqueConstraintIdentification) {
			messagesParameters := []interface{}{client.Identification}
			return &errors.GracefulError{ErrorCode: errors.IdentificationDuplicate, MessagesParameters: messagesParameters}
		}
	}
	return error
}

func (client *Client) Update() error {
	if exist, error := clientExist(client.ID); !exist {
		return error
	}
	error := config.DB.Save(client).Error
	if error != nil {
		if postgres.IsUniqueConstraintError(error, UniqueConstraintIdentification) {
			messagesParameters := []interface{}{client.Identification}
			return &errors.GracefulError{ErrorCode: errors.IdentificationDuplicate, MessagesParameters: messagesParameters}
		}
		return error
	}
	client.Addresses, error = findAddressesByClientId(client.ID)
	return error
}

func FindClientByID(clientID uint) (*Client, error) {
	var client Client
	response := config.DB.First(&client, clientID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{clientID}
			return nil, &errors.RecordNotFound{ErrorCode: errors.ClientNotExist, MessagesParameters: messagesParameters}
		}
		return nil, error
	}
	return &client, nil
}

func clientExist(clientID uint) (bool, error) {
	if _, error := FindClientByID(clientID); error != nil {
		if _, ok := error.(*errors.RecordNotFound); ok {
			return false, error
		}
		return false, error
	}
	return true, nil
}

func validateClientAddress(addresses []Address) error {
	if len(addresses) <= 0 {
		return &errors.GracefulError{ErrorCode: errors.AddressRequired}
	}
	for _, address := range addresses {
		if _, error := findCityByID(address.CityID); error != nil {
			return error
		}
	}
	return nil
}
