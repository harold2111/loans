package domain

// Repository provides access a client store.
type ClientRepository interface {
	FindAll() ([]Client, error)
	Find(clientID uint) (Client, error)
	Create(client *Client) error
	Update(client *Client) error
	Delete(client *Client) error
	FindAddressesByClientID(addressID uint) ([]Address, error)
	FindAddressByIDAndClientID(addressID uint, ClientID uint) (*Address, error)
	CreateAddressClient(address *Address) error
	UpdateAddressClient(address *Address) error
	DeleteAddressClient(address *Address) error
}
