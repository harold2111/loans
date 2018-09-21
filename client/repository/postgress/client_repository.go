package postgres

import (
	"loans/client"
	"loans/errors"
	"loans/models"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type clientRepository struct {
	db *gorm.DB
}

const (
	uniqueConstraintIdentification = "uix_clients_identification"
)

// NewClientRepository returns a new instance of a Postgres client repository.
func NewClientRepository(db *gorm.DB) (client.ClientRepository, error) {
	r := &clientRepository{
		db: db,
	}
	return r, nil
}

func (r *clientRepository) FindAll() ([]models.Client, error) {
	var clients []models.Client
	response := r.db.Preload("Addresses").Find(&clients)
	if error := response.Error; error != nil {
		return nil, error
	}
	return clients, nil
}

func (r *clientRepository) Find(clientID uint) (models.Client, error) {
	var client models.Client
	response := r.db.Preload("Addresses").First(&client, clientID)
	if error := response.Error; error != nil {
		if response.RecordNotFound() {
			messagesParameters := []interface{}{clientID}
			return client, &errors.RecordNotFound{ErrorCode: errors.ClientNotExist, MessagesParameters: messagesParameters}
		}
		return client, error
	}
	return client, nil
}

func (r *clientRepository) Store(client *models.Client) error {
	removeIDs(client)
	error := r.db.Create(client).Error
	if error != nil {
		if isUniqueConstraintError(error, uniqueConstraintIdentification) {
			messagesParameters := []interface{}{client.Identification}
			return &errors.GracefulError{ErrorCode: errors.IdentificationDuplicate, MessagesParameters: messagesParameters}
		}
	}
	return error
}

func (r *clientRepository) Update(client *models.Client) error {
	if exist, err := r.ClientExist(client.ID); !exist {
		return err
	}
	currentAddresses, err := r.FindClientAddress(client.ID)
	if err != nil {
		return err
	}
	toCreateAddresses, toUpdateAddresses, err := addressesToCreateUpdate(currentAddresses, client.Addresses)
	if err != nil {
		return err
	}
	toCreateAddreses := 
	err = r.db.Save(client).Error
	if err != nil {
		if isUniqueConstraintError(err, uniqueConstraintIdentification) {
			messagesParameters := []interface{}{client.Identification}
			return &errors.GracefulError{ErrorCode: errors.IdentificationDuplicate, MessagesParameters: messagesParameters}
		}
	}
	return err
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

func (r *clientRepository) FindClientAddress(clientID uint) ([]models.Address, error) {
	var addresses []models.Address
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

func removeIDs(client *models.Client) {
	client.ID = 0
	for index := 0; index < len(client.Addresses); index++ {
		client.Addresses[index].ID = 0
	}
}

func addressesToCreateUpdate(currents []models.Address, addresses []models.Address) ([]models.Address, []models.Address, error) {
	var addressesToCreate []models.Address
	var addressesToUpdate []models.Address
	
	for _, address := range addresses {
		if address.ID == 0 {
			addressesToCreate = append(addressesToCreate, address)
		}else {
			if replaceUpdateAddressInCurrentAddress(currents, address.ID){
				addressesToUpdate = append(addressesToUpdate, address)
			}else{
				messagesParameters := []interface{}{address.ID}
			    return nil, nil, &errors.GracefulError{
					ErrorCode: errors.ClientNotAddressFound, 
					MessagesParameters: messagesParameters,
				}
			}
			
		}
	}
	return addressesToCreate, addressesToUpdate, nil
}

func replaceUpdateAddressInCurrentAddress(currentAddresses []models.Address, toUpdateAddress models.Address) error {
	for i, address := range addresses {
		if address.ID == toFindAddressID {
			currentAddresses[i] = toUpdateAddress
		}
	}
	return messagesParameters := []interface{}{toUpdateAddress.ID}
	return &errors.GracefulError{
		ErrorCode: errors.ClientNotAddressFound, 
		MessagesParameters: messagesParameters,
	}
}

func isUniqueConstraintError(err error, constraintName string) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" && pqErr.Constraint == constraintName
	}
	return false
}
