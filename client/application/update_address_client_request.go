package application

type UpdateAddressClientRequest struct {
	ID            uint   `json:"id"`
	StreetAddress string `json:"streetAddress"`
	ClientID      uint   `json:"clientID"`
	DepartmentID  uint   `json:"departmentID"`
	CityID        uint   `json:"cityID"`
}
