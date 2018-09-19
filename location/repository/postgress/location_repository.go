package postgres

import (
	"loans/errors"
	"loans/location"
	"loans/models"

	"github.com/jinzhu/gorm"
)

type locationRepository struct {
	db *gorm.DB
}

// NewLocationRepositoryy returns a new instance of a Postgres location repository.
func NewLocationRepositoryy(db *gorm.DB) (location.LocationRepository, error) {
	r := &locationRepository{
		db: db,
	}
	return r, nil
}

func (r *locationRepository) FindAllDepartments() ([]models.Department, error) {
	var departments []models.Department
	response := r.db.Find(&departments)
	if error := response.Error; error != nil {
		return nil, error
	}
	return departments, nil
}

func (r *locationRepository) FindCity(cityID uint) (*models.City, error) {
	var city models.City
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

func (r *locationRepository) FindCitiesByDepartmentID(departmentID uint) ([]models.City, error) {
	var cities []models.City
	response := r.db.Find(&cities, "department_id = ?", departmentID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{departmentID}
			return nil, &errors.RecordNotFound{ErrorCode: errors.NotCitiesForDepartment,
				MessagesParameters: messagesParameters}
		}
		return nil, error
	}
	return cities, nil
}
