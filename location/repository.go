package location

// Repository provides access a address store.
type Repository interface {
	FindAllDepartments() ([]Department, error)
	FindCitiesByDepartmentID(departmentID uint) ([]City, error)
	FindCity(cityID uint) (*City, error)
}
