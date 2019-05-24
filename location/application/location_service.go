package application

import (
	locationDomain "loans/location/domain"
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

func (s *LocationService) FindAllDepartments() ([]locationDomain.Department, error) {
	return s.locationRepository.FindAllDepartments()
}

func (s *LocationService) FindCitiesByDepartmentID(departmentID uint) ([]locationDomain.City, error) {
	return s.locationRepository.FindCitiesByDepartmentID(departmentID)
}
