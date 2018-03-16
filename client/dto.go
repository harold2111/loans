package client

type ClientCommonFields struct {
	Identification string `json:"identification" validate:"required"`
	FirstName      string `json:"firstName" validate:"required"`
	LastName       string `json:"lastName" validate:"required"`
	Telephone1     string `json:"telephone1" validate:"required"`
	Telephone2     string `json:"telephone2,omitempty"`
}

// Request
type CreateClientRequest struct {
	ClientCommonFields
	Address CreateAddress `json:"address" validate:"required"`
}

type CreateAddress struct {
	Address string `json:"address" validate:"required"`
	CityID  uint   `json:"cityID" validate:"required"`
}

//Response
type ClientResponse struct {
	ID uint `json:"id"`
	ClientCommonFields
	Address AddressResponse `json:"address" validate:"required"`
}

type AddressResponse struct {
	ID       uint `json:"id"`
	ClientID uint `json:"clientID"`
	CreateAddress
}

type UpdateClientRequest struct {
	ClientCommonFields
}
