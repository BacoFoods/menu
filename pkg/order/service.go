package order

import (
	"fmt"
	products "github.com/BacoFoods/menu/pkg/product"
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
	product    products.Repository
}

func NewService(repository Repository, table tables.Repository, product products.Repository) service {
	return service{repository, table, product}
}

func (s service) Create(order *Order) (*Order, error) {
	productIDs := order.GetProductIDs()
	prods, err := s.product.GetByIDs(productIDs)
	if err != nil {
		shared.LogError("error getting products", LogService, "Create", err, productIDs)
		return nil, fmt.Errorf(ErrorOrderCreation)
	}

	order.SetItems(prods)

	newOrder, err := s.repository.Create(order)
	if err != nil {
		shared.LogError("error creating order", LogService, "Create", err, *order)
		return nil, fmt.Errorf(ErrorOrderCreation)
	}

	if _, err := s.table.SetOrder(newOrder.TableID, &newOrder.ID); err != nil {
		return nil, err
	}

	orderDB, err := s.repository.Get(fmt.Sprintf("%d", newOrder.ID))
	if err != nil {
		shared.LogError("error getting order", LogService, "Create", err, newOrder.ID)
		return nil, fmt.Errorf(ErrorOrderCreation)
	}

	return orderDB, nil
}
