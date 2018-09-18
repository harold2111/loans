package postgres

import (
	"loans/client"
	"loans/errors"

	"github.com/jinzhu/gorm"
)

type clientRepository struct {
	db *gorm.DB
}

const (
	uniqueConstraintIdentification = "uix_clients_identification"
)

// NewClientRepository returns a new instance of a Postgres client repository.
func NewClientRepository(db *gorm.DB) (client.Repository, error) {
	r := &clientRepository{
		db: db,
	}
	return r, nil
}

func (r *clientRepository) FindAll() ([]client.Client, error) {
	var clients []client.Client
	response := r.db.Find(&clients)
	if error := response.Error; error != nil {
		return nil, error
	}
	return clients, nil
}

func (r *clientRepository) Find(clientID uint) (client.Client, error) {
	var client client.Client
	response := r.db.First(&client, clientID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{clientID}
			return client, &errors.RecordNotFound{ErrorCode: errors.ClientNotExist, MessagesParameters: messagesParameters}
		}
		return client, error
	}
	return client, nil
}

func (r *clientRepository) Store(client *client.Client) error {
	error := r.db.Create(client).Error
	if error != nil {
		if isUniqueConstraintError(error, uniqueConstraintIdentification) {
			messagesParameters := []interface{}{client.Identification}
			return &errors.GracefulError{ErrorCode: errors.IdentificationDuplicate, MessagesParameters: messagesParameters}
		}
	}
	return error
}

func (r *clientRepository) Update(client *client.Client) error {
	error := r.db.Save(client).Error
	if error != nil {
		if isUniqueConstraintError(error, uniqueConstraintIdentification) {
			messagesParameters := []interface{}{client.Identification}
			return &errors.GracefulError{ErrorCode: errors.IdentificationDuplicate, MessagesParameters: messagesParameters}
		}
	}
	return error
}

func (r *clientRepository) ClientExist(clientID uint) (bool, error) {
	if _, error := r.Find(clientID); error != nil {
		if _, ok := error.(*errors.RecordNotFound); ok {
			return false, error
		}
		return false, error
	}
	return true, nil
}

func (r *clientRepository) FindClientAddress(clientID uint) ([]client.Address, error) {
	var addresses []client.Address
	response := r.db.Find(&addresses, "client_id = ?", clientID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{clientID}
			return nil, &errors.RecordNotFound{ErrorCode: errors.ClientNotAddressFound, MessagesParameters: messagesParameters}
		}
		return nil, error
	}
	return addresses, nil
}
