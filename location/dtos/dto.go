package dtos

type DepartmentResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CityResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
