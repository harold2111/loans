package domain

import "loans/shared/models"

type LocationRepository interface {
	FindAllDepartments() ([]models.Department, error)
	FindCitiesByDepartmentID(departmentID uint) ([]models.City, error)
	FindCity(cityID uint) (*models.City, error)
}
