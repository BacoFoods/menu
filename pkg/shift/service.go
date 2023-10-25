package shift

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/account"
	"github.com/BacoFoods/menu/pkg/shared"
	"time"
)

const (
	LogService string = "pkg/shift/service"
)

type Service interface {
	Open(accountUUID string, startBalance float32) (*Shift, error)
	Close(accountUUID string, endBalance float32) (*Shift, error)
}

type service struct {
	repository        Repository
	accountRepository account.Repository
}

func NewService(repository Repository, accountRepository account.Repository) service {
	return service{repository, accountRepository}
}

func (s service) Open(accountUUID string, startBalance float32) (*Shift, error) {
	acc, err := s.accountRepository.GetByUUID(accountUUID)
	if err != nil {
		shared.LogError("failed to get account", LogService, "Open", err)
		return nil, fmt.Errorf(ErrorGettingAccount)
	}

	now := time.Now()

	shift := &Shift{
		StoreID:      acc.StoreID,
		BrandID:      acc.BrandID,
		AccountID:    &acc.Id,
		StartTime:    &now,
		EndTime:      nil,
		StartBalance: startBalance,
		EndBalance:   0,
	}
	return s.repository.Create(shift)
}

func (s service) Close(accountUUID string, endBalance float32) (*Shift, error) {
	acc, err := s.accountRepository.GetByUUID(accountUUID)
	if err != nil {
		shared.LogError("failed to get account", LogService, "Close", err)
		return nil, fmt.Errorf(ErrorGettingAccount)
	}

	openShift, err := s.repository.GetOpenShift(acc.StoreID)
	if err != nil {
		shared.LogError("failed to get open shift", LogService, "Close", err)
		return nil, fmt.Errorf(ErrorGettingOpenShift)
	}

	now := time.Now()
	openShift.EndTime = &now
	openShift.EndBalance = endBalance

	return s.repository.Update(openShift)
}
