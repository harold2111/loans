package application

import clientDomain "loans/client/domain"

type GetClientResponse struct {
	ID             uint                       `json:"id"`
	Identification string                     `json:"identification"`
	FirstName      string                     `json:"firstName"`
	LastName       string                     `json:"lastName"`
	Telephone1     string                     `json:"telephone1"`
	Telephone2     string                     `json:"telephone2,omitEmpty"`
	Email          string                     `json:"email"`
	Addresses      []GetAddressClientResponse `json:"addresses,omitempty"`
}

func (g *GetClientResponse) fillFromClient(client clientDomain.Client) {
	g.ID = client.ID
	g.Identification = client.Identification
	g.FirstName = client.FirstName
	g.LastName = client.LastName
	g.Telephone1 = client.Telephone1
	g.Telephone2 = client.Telephone2
	g.Email = client.Email
	if len(client.Addresses) > 0 {
		for _, address := range client.Addresses {
			var getAddressClientResponse GetAddressClientResponse
			getAddressClientResponse.fillFromAddress(address)
			g.Addresses = append(g.Addresses, getAddressClientResponse)
		}
	}
}
