package domain

type LocationRepository interface {
	FindAllDepartments() ([]Department, error)
	FindCitiesByDepartmentID(departmentID uint) ([]City, error)
	FindCity(cityID uint) (*City, error)
}
