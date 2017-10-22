package dtos

type CreateAddress struct {
	Address string `json:"address" validate:"required"`
	CityID  uint   `json:"cityID" validate:"required"`
}

type AddressResponse struct {
	ID       uint `json:"id"`
	ClientID uint `json:"clientID"`
	CreateAddress
}
