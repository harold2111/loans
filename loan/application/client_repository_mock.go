package application

import "github.com/harold2111/loans/client/domain"

type clientRepositoryMock struct {
}

func (r *clientRepositoryMock) Create(client *domain.Client) error {
	panic("not implemented") // TODO: Implement
}

func (r *clientRepositoryMock) FindAll() ([]domain.Client, error) {
	panic("not implemented") // TODO: Implement
}

func (r *clientRepositoryMock) Find(clientID uint) (domain.Client, error) {
	panic("not implemented") // TODO: Implement
}

func (r *clientRepositoryMock) Update(client *domain.Client) error {
	panic("not implemented") // TODO: Implement
}

func (r *clientRepositoryMock) Delete(client *domain.Client) error {
	panic("not implemented") // TODO: Implement
}
