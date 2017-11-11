package postgres

import (
	"loans/errors"
	"loans/location"

	"github.com/jinzhu/gorm"
)

type locationRepository struct {
	db *gorm.DB
}

// NewLocationRepositoryy returns a new instance of a Postgres location repository.
func NewLocationRepositoryy(db *gorm.DB) (location.Repository, error) {
	r := &locationRepository{
		db: db,
	}
	return r, nil
}

func (r *locationRepository) FindCity(cityID uint) (*location.City, error) {
	var city location.City
	response := r.db.First(&city, cityID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{cityID}
			return nil, &errors.RecordNotFound{ErrorCode: errors.CityNotExist, MessagesParameters: messagesParameters}
		}
		return nil, error
	}
	return &city, nil
}
