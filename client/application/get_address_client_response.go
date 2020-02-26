package application

import clientDomain "github.com/harold2111/loans/client/domain"

type GetAddressClientResponse struct {
	ID            uint   `json:"id"`
	StreetAddress string `json:"streetAddress"`
	ClientID      uint   `json:"clientID"`
	DepartmentID  uint   `json:"departmentID"`
	CityID        uint   `json:"cityID"`
}

func (g *GetAddressClientResponse) fillFromAddress(client clientDomain.Address) {
	g.ID = client.ID
	g.StreetAddress = client.StreetAddress
	g.ClientID = client.ClientID
	g.DepartmentID = client.DepartmentID
	g.CityID = client.CityID
}
