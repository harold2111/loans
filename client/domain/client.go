package domain

import (
	"github.com/harold2111/loans/shared/errors"
	"github.com/harold2111/loans/shared/utils"
	"time"
)

type Client struct {
	ID             uint `gorm:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `sql:"index"`
	Identification string     `gorm:"not null; unique_index"`
	FirstName      string     `gorm:"not null"`
	LastName       string     `gorm:"not null"`
	Telephone1     string     `gorm:"not null"`
	Telephone2     string
	Email          string
	Addresses      []Address
}

func NewClientForCreate(
	identification string,
	firstName string,
	lastName string,
	telephone1 string,
	telephone2 string,
	email string,
	addresses []Address) (Client, error) {

	client := Client{
		Identification: identification,
		FirstName:      firstName,
		LastName:       lastName,
		Telephone1:     telephone1,
		Telephone2:     telephone2,
		Email:          email,
		Addresses:      addresses,
	}
	if error := client.validateForCreation(); error != nil {
		return Client{}, error
	}
	return client, nil
}

func NewClientForUpdate(
	id uint,
	identification string,
	firstName string,
	lastName string,
	telephone1 string,
	telephone2 string,
	email string,
	addresses []Address) (Client, error) {

	client := Client{
		ID:             id,
		Identification: identification,
		FirstName:      firstName,
		LastName:       lastName,
		Telephone1:     telephone1,
		Telephone2:     telephone2,
		Email:          email,
		Addresses:      addresses,
	}
	if error := client.validateForCreation(); error != nil {
		return Client{}, error
	} else if error := utils.ValidateVar("id", client.ID, "required"); error != nil {
		return Client{}, error
	}
	return client, nil
}

func (c *Client) validateForCreation() error {
	if len(c.Addresses) == 0 {
		return &errors.GracefulError{ErrorCode: errors.AtLeastOneAddress}
	} else if error := utils.ValidateVar("identification", c.Identification, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("firstName", c.FirstName, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("lastName", c.LastName, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("telephone1", c.Telephone1, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("email", c.Email, "required,email"); error != nil {
		return error
	}
	return nil
}
