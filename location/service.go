package location

import "loans/models"

// Service is the interface that provides deparment methods.
type LocationService interface {
	FindAllDepartments() ([]models.Department, error)
	FindCitiesByDepartmentID(departmentID uint) ([]models.City, error)
}
