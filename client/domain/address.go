package domain

import (
	"github.com/harold2111/loans/shared/utils"
	"time"
)

type Address struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	StretAddress string `gorm:"not null"`
	ClientID     uint   `gorm:"not null"`
	DepartmentID uint   `gorm:"not null"`
	CityID       uint   `gorm:"not null"`
}

func NewAddessForCreateClient(
	stretAddress string,
	departmentID uint,
	cityID uint) (Address, error) {
	address := Address{
		StretAddress: stretAddress,
		DepartmentID: departmentID,
		CityID:       cityID,
	}
	if error := address.validateForCreationOfNewClient(); error != nil {
		return Address{}, error
	}
	return address, nil
}

func NewAddessForUpdateClient(
	id uint,
	stretAddress string,
	departmentID uint,
	cityID uint) (Address, error) {
	address := Address{
		ID:           id,
		StretAddress: stretAddress,
		DepartmentID: departmentID,
		CityID:       cityID,
	}
	if error := address.validateForCreationOfNewClient(); error != nil {
		return Address{}, error
	}
	return address, nil
}

func (a *Address) validateForCreationOfNewClient() error {
	if error := utils.ValidateVar("stretAddress", a.StretAddress, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("deparmentID", a.DepartmentID, "required"); error != nil {
		return error
	} else if error := utils.ValidateVar("cityID", a.CityID, "required"); error != nil {
		return error
	}
	return nil
}
