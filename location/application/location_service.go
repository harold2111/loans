package application

import (
	locationDomain "loans/location/domain"
	"loans/shared/models"
)

type LocationService struct {
	locationRepository locationDomain.LocationRepository
}

// NewLocationService creates a deparment service with necessary dependencies.
func NewLocationService(locationRepository locationDomain.LocationRepository) LocationService {
	return LocationService{
		locationRepository: locationRepository,
	}
}

func (s *LocationService) FindAllDepartments() ([]models.Department, error) {
	return s.locationRepository.FindAllDepartments()
}

func (s *LocationService) FindCitiesByDepartmentID(departmentID uint) ([]models.City, error) {
	return s.locationRepository.FindCitiesByDepartmentID(departmentID)
}
