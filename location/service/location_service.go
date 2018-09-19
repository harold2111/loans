package location

import (
	"loans/location"
	"loans/models"
)

type locationService struct {
	locationRepository location.LocationRepository
}

// NewLocationService creates a deparment service with necessary dependencies.
func NewLocationService(locationRepository location.LocationRepository) location.LocationService {
	return &locationService{
		locationRepository: locationRepository,
	}
}

func (s *locationService) FindAllDepartments() ([]models.Department, error) {
	return s.locationRepository.FindAllDepartments()
}

func (s *locationService) FindCitiesByDepartmentID(departmentID uint) ([]models.City, error) {
	return s.locationRepository.FindCitiesByDepartmentID(departmentID)
}
