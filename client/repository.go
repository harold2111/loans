package client

// Repository provides access a client store.
type Repository interface {
	Store(client *Client) error
	Update(client *Client) error
	Find(clientID uint) (*Client, error)
	ClientExist(clientID uint) (bool, error)
	StoreClientAddresses(clientID uint, addresses *[]Address) error
	FindClientAddress(clientID uint) ([]Address, error)
}
