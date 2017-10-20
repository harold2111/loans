package models

import (
	"loans/config"
	"loans/errors"

	"github.com/jinzhu/gorm"
)

type Address struct {
	gorm.Model
	Address  string `gorm:"not null"`
	Client   Client
	ClientID uint `gorm:"not null"`
	City     City
	CityID   uint `gorm:"not null"`
}

func findAddressesByClientId(clientID uint) ([]Address, error) {
	var addresses []Address
	response := config.DB.Find(&addresses, "client_id = ?", clientID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{clientID}
			return nil, &errors.RecordNotFound{ErrorCode: errors.ClientNotAddressFound, MessagesParameters: messagesParameters}
		}
		return nil, error
	}
	return addresses, nil
}
