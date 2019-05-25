package application

type GetAddressClientResponse struct {
	ID           uint   `json:"id"`
	StretAddress string `json:"stretAddress"`
	ClientID     uint   `json:"clientID"`
	DepartmentID uint   `json:"departmentID"`
	CityID       uint   `json:"cityID"`
}
