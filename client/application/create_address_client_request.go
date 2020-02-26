package application

type CreateAddressClientRequest struct {
	StreetAddress string `json:"streetAddress"`
	ClientID      uint   `json:"clientID"`
	DepartmentID  uint   `json:"departmentID"`
	CityID        uint   `json:"cityID"`
}
