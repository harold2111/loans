package models

import (
	"loans/config"
	"loans/errors"

	"github.com/jinzhu/gorm"
)

type City struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Department   Department
	DepartmentID uint `gorm:"not null"`
}

func findCityByID(cityID uint) (*City, error) {
	var city City
	response := config.DB.First(&city, cityID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{cityID}
			return nil, &errors.RecordNotFound{ErrorCode: errors.CityNotExist, MessagesParameters: messagesParameters}
		}
		return nil, error
	}
	return &city, nil
}
