package order

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/BacoFoods/menu/pkg/tables"
)

const (
	LogService string = "pkg/order/service"
)

type Service interface {
	Create(*Order) (*Order, error)
}

type service struct {
	repository Repository
	table      tables.Repository
}

func NewService(repository Repository, table tables.Repository) service {
	return service{repository, table}
}

func (s service) Create(order *Order) (*Order, error) {
	newOrder, err := s.repository.Create(order)
	if err != nil {
		shared.LogError("error creating order", LogService, "Create", err, *order)
		return nil, err
	}

	if _, err := s.table.SetOrder(newOrder.TableID, &newOrder.ID); err != nil {
		return nil, err
	}

	orderDB, err := s.repository.Get(fmt.Sprintf("%d", newOrder.ID))
	if err != nil {
		shared.LogError("error getting order", LogService, "Create", err, newOrder.ID)
		return nil, err
	}

	return orderDB, nil
}
