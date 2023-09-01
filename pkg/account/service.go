package account

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
)

const (
	LogService string = "pkg/account/service"
)

type Service interface {
	Create(*Account) (*Account, error)
	CreatePinUser(*Account) (*Account, error)
	Login(username, password string) (*Account, error)
	LoginPin(pin int) (*Account, error)
	Delete(id string) error
	Find(filter map[string]any) ([]Account, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

func (s service) Create(account *Account) (*Account, error) {
	if err := account.HashPassword(); err != nil {
		shared.LogError(ErrorAccountCreation, LogService, "Create", err, account)
		return nil, fmt.Errorf(ErrorAccountCreation)
	}

	return s.repository.Create(account)
}

func (s service) CreatePinUser(account *Account) (*Account, error) {
	account.HashPin()

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

func (s service) LoginPin(pin int) (*Account, error) {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%d", pin)))
	hashBytes := hasher.Sum(nil)
	hashed := hex.EncodeToString(hashBytes)

	account, err := s.repository.Get(hashed)
	if err != nil {
		shared.LogError(ErrorAccountLogin, LogService, "LoginPin", err, pin)
		return nil, err
	}

	return account, nil
}

func (s service) Delete(id string) error {
	return s.repository.Delete(id)
}

func (s service) Find(filter map[string]any) ([]Account, error) {
	return s.repository.Find(filter)
}
