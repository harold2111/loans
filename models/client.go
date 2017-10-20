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

func (client *Client) Save() error {
	if error := validateClientAddress(client.Addresses); error != nil {
		fmt.Println(error)
		return error
	}
	error := config.DB.Create(client).Error
	if error != nil {
		fmt.Println(error)
		if postgres.IsUniqueConstraintError(error, UniqueConstraintIdentification) {
			messagesParameters := []interface{}{client.Identification}
			return &errors.GracefulError{ErrorCode: errors.IdentificationDuplicate, MessagesParameters: messagesParameters}
		}
	}
	return error
}

func (client *Client) UpdateClient() error {
	error := config.DB.Save(client).Error
	if error != nil {
		fmt.Println(error)
		if postgres.IsUniqueConstraintError(error, UniqueConstraintIdentification) {
			messagesParameters := []interface{}{client.Identification}
			return &errors.GracefulError{ErrorCode: errors.IdentificationDuplicate, MessagesParameters: messagesParameters}
		}
		return error
	}
	client.Addresses, error = findAddressesByClientId(client.ID)
	return nil
}

func validateClientAddress(addresses []Address) error {
	if len(addresses) <= 0 {
		return &errors.GracefulError{ErrorCode: errors.AddressRequired}
	}
	for _, address := range addresses {
		if _, error := findCityByID(address.CityID); error != nil {
			if recordNotFound, ok := error.(*errors.RecordNotFound); ok {
				return &errors.GracefulError{ErrorCode: recordNotFound.ErrorCode, MessagesParameters: recordNotFound.MessagesParameters}
			}
			return error
		}
	}
	return nil
}
