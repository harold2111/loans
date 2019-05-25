package application

type CreateAddressClientRequest struct {
	StretAddress string `json:"stretAddress"`
	ClientID     uint   `json:"clientID"`
	DepartmentID uint   `json:"departmentID"`
	CityID       uint   `json:"cityID"`
}
