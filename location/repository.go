package location

import "loans/models"

type LocationRepository interface {
	FindAllDepartments() ([]models.Department, error)
	FindCitiesByDepartmentID(departmentID uint) ([]models.City, error)
	FindCity(cityID uint) (*models.City, error)
}
