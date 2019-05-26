package domain

// Repository provides access a client store.
type ClientRepository interface {
	Create(client *Client) error
	FindAll() ([]Client, error)
	Find(clientID uint) (Client, error)
	Update(client *Client) error
	Delete(client *Client) error
}
