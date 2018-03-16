package location

// Repository provides access a address store.
type Repository interface {
	FindCity(cityID uint) (*City, error)
	FindCitiesByDepartment(departmentID uint) ([]City, error)
	FindAllDepartments() ([]Department, error)
}
