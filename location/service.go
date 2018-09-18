package location

type service struct {
	locationRepository Repository
}

// Service is the interface that provides deparment methods.
type Service interface {
	FindAllDepartments() ([]Department, error)
	FindCitiesByDepartmentID(departmentID uint) ([]City, error)
}

// NewService creates a deparment service with necessary dependencies.
func NewService(locationRepository Repository) Service {
	return &service{
		locationRepository: locationRepository,
	}
}

func (s *service) FindAllDepartments() ([]Department, error) {
	return s.locationRepository.FindAllDepartments()
}

func (s *service) FindCitiesByDepartmentID(departmentID uint) ([]City, error) {
	return s.locationRepository.FindCitiesByDepartmentID(departmentID)
}
