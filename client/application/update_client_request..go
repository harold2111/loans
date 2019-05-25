package application

type UpdateClientRequest struct {
	ID             uint                         `json:"id"`
	Identification string                       `json:"identification"`
	FirstName      string                       `json:"firstName"`
	LastName       string                       `json:"lastName"`
	Telephone1     string                       `json:"telephone1"`
	Telephone2     string                       `json:"telephone2,omitEmpty"`
	Email          string                       `json:"email"`
	Addresses      []UpdateAddressClientRequest `json:"addresses,omitempty"`
}
