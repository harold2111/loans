package domain

import (
	"time"

	"github.com/harold2111/loans/shared/utils"
)

type Address struct {
	ID            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StreetAddress string `gorm:"not null"`
	ClientID      uint   `gorm:"not null"`
	DepartmentID  uint   `gorm:"not null"`
	CityID        uint   `gorm:"not null"`
}

func NewAddessForCreateClient(
	streetAddress string,
	departmentID uint,
	cityID uint) (Address, error) {
	address := Address{
		StreetAddress: streetAddress,
		DepartmentID:  departmentID,
		CityID:        cityID,
	}
	if error := address.validateForCreationOfNewClient(); error != nil {
		return Address{}, error
	}
	return address, nil
}

func NewAddessForUpdateClient(
	id uint,
	streetAddress string,
	departmentID uint,
	cityID uint) (Address, error) {
	address := Address{
		ID:            id,
		StreetAddress: streetAddress,
		DepartmentID:  departmentID,
		CityID:        cityID,
	}
	if error := address.validateForCreationOfNewClient(); error != nil {
		return Address{}, error
	}
	return address, nil
}

func (a *Address) validateForCreationOfNewClient() error {
	if error := utils.ValidateVar("streetAddress", a.StreetAddress, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("departmentID", a.DepartmentID, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("cityID", a.CityID, "required"); error != nil {
		return error
	}
	return nil
}
