package dtos

type ClientCommonData struct {
	Identification string `json:"identification" validate:"required"`
	FirstName      string `json:"firstName" validate:"required"`
	LastName       string `json:"lastName" validate:"required"`
	Telephone1     string `json:"telephone1" validate:"required"`
	Telephone2     string `json:"telephone2,omitempty"`
}

type CreateClient struct {
	ClientCommonData
	Addresses []CreateAddress `json:"addresses" validate:"required"`
}

type UpdateClient struct {
	ClientCommonData
}

type CreateAddress struct {
	Address string `json:"address" validate:"required"`
	CityID  uint   `json:"cityID" validate:"required"`
}

type AddressDTO struct {
	ID uint `json:"id"`
	CreateAddress
}

type ClientDTO struct {
	ID uint `json:"id"`
	ClientCommonData
	Addresses []AddressDTO `json:"addresses,omitempty" validate:"required"`
}
