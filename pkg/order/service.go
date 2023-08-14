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
	UpdateTable(orderID, tableID uint64) (*Order, error)
	Get(string) (*Order, error)
	Find(filter map[string]any) ([]Order, error)
	UpdateSeats(orderID string, seats int) (*Order, error)
	AddProduct(orderID, productID string) (*Order, error)
	RemoveProduct(orderID, productID string) (*Order, error)
	UpdateProduct(product *OrderItem) (*Order, error)

	CreateOrderType(orderType *OrderType) (*OrderType, error)
	FindOrderType(filter map[string]any) ([]OrderType, error)
	GetOrderType(orderTypeID string) (*OrderType, error)
	UpdateOrderType(orderTypeID string, orderType *OrderType) (*OrderType, error)
	DeleteOrderType(orderTypeID string) error
}

type service struct {
	repository Repository
	table      tables.Repository
	product    products.Repository
}

func NewService(repository Repository, table tables.Repository, product products.Repository) service {
	return service{repository, table, product}
}

// Orders

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

func (s service) UpdateTable(orderID, tableID uint64) (*Order, error) {
	order, err := s.repository.Get(fmt.Sprintf("%d", orderID))
	if err != nil {
		shared.LogError("error getting order", LogService, "UpdateTable", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	oldTableID := *order.TableID
	newTableID := uint(tableID)

	if oldTableID == newTableID {
		return order, nil
	}

	if _, err := s.table.SetOrder(&newTableID, &order.ID); err != nil {
		return nil, err
	}

	if _, err := s.table.RemoveOrder(&oldTableID); err != nil {
		return nil, err
	}

	order.TableID = &newTableID
	orderDB, err := s.repository.Update(order)
	if err != nil {
		shared.LogError("error updating order", LogService, "UpdateTable", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	return orderDB, nil
}

func (s service) Get(id string) (*Order, error) {
	return s.repository.Get(id)
}

func (s service) Find(filter map[string]any) ([]Order, error) {
	return s.repository.Find(filter)
}

func (s service) UpdateSeats(orderID string, seats int) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "UpdateSeats", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	if order.Seats == seats {
		return order, nil
	}

	order.Seats = seats
	orderDB, err := s.repository.Update(order)
	if err != nil {
		shared.LogError("error updating order", LogService, "UpdateSeats", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	return orderDB, nil
}

func (s service) AddProduct(orderID, productID string) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "AddProduct", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	product, err := s.product.Get(productID)
	if err != nil {
		shared.LogError("error getting product", LogService, "AddProduct", err, productID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	order.AddProduct(product)

	orderDB, err := s.repository.Update(order)
	if err != nil {
		shared.LogError("error updating order", LogService, "AddProduct", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	return orderDB, nil
}

func (s service) RemoveProduct(orderID, productID string) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "RemoveProduct", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	product, err := s.product.Get(productID)
	if err != nil {
		shared.LogError("error getting product", LogService, "RemoveProduct", err, productID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	order.RemoveProduct(product)

	orderDB, err := s.repository.Update(order)
	if err != nil {
		shared.LogError("error updating order", LogService, "RemoveProduct", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	return orderDB, nil
}

func (s service) UpdateProduct(product *OrderItem) (*Order, error) {
	orderItem, err := s.repository.UpdateOrderItem(product)
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderItemUpdate)
	}

	order, err := s.repository.Get(fmt.Sprintf("%d", orderItem.OrderID))
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	return order, nil
}

// Order Types

func (s service) CreateOrderType(orderType *OrderType) (*OrderType, error) {
	return s.repository.CreateOrderType(orderType)
}

func (s service) FindOrderType(filter map[string]any) ([]OrderType, error) {
	return s.repository.FindOrderType(filter)
}

func (s service) GetOrderType(orderTypeID string) (*OrderType, error) {
	return s.repository.GetOrderType(orderTypeID)
}

func (s service) UpdateOrderType(orderTypeID string, orderType *OrderType) (*OrderType, error) {
	return s.repository.UpdateOrderType(orderTypeID, orderType)
}

func (s service) DeleteOrderType(orderTypeID string) error {
	return s.repository.DeleteOrderType(orderTypeID)
}
