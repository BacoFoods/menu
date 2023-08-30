package account

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
)

const (
	LogService string = "pkg/account/service"
)

type Service interface {
	Create(*Account) (*Account, error)
	Login(username, password string) (*Account, error)
	Delete(id string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Create(account *Account) (*Account, error) {

	return s.repository.Create(account)
}

func (s service) Login(username, password string) (*Account, error) {
	account, err := s.repository.Get(username)
	if err != nil {
		shared.LogError(ErrorAccountLogin, LogService, "Login", err, username)
		return nil, err
	}

	if !account.CheckPassword(password) {
		return nil, fmt.Errorf(ErrorAccountInvalidPassword)
	}

	return account, nil
}

func (s service) Delete(id string) error {
	return s.repository.Delete(id)
}
