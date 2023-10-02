package order

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	LogDBRepository string = "pkg/order/db_repository"
	StatusActive    string = "active"
	StatusCreate    string = "create"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

// Order methods

// Create method for create a new order in database
func (r *DBRepository) Create(order *Order) (*Order, error) {
	if err := r.db.Save(order).Error; err != nil {
		shared.LogError("error creating order", LogDBRepository, "Create", err, *order)
		return nil, err
	}

	return order, nil
}

// Get method for get an order from database
func (r *DBRepository) Get(orderID string) (*Order, error) {
	if orderID == "" {
		shared.LogWarn("error getting order", LogDBRepository, "Get", shared.ErrorIDEmpty)
		return nil, shared.ErrorIDEmpty
	}

	var order Order
	if err := r.db.
		Preload(clause.Associations).
		Preload("Items.Modifiers").
		First(&order, orderID).Error; err != nil {
		shared.LogError("error getting order", LogDBRepository, "Get", err, orderID)
		return nil, err
	}

	return &order, nil
}

// Update method for update an order in database
func (r *DBRepository) Update(order *Order) (*Order, error) {
	if err := r.db.Save(order).Error; err != nil {
		shared.LogError("error updating order", LogDBRepository, "Update", err, *order)
		return nil, err
	}

	return order, nil
}

// Find method for find orders in database
func (r *DBRepository) Find(filter map[string]any) ([]Order, error) {
	tx := r.db.
		Preload(clause.Associations).
		Preload("Items.Modifiers")

	if status, ok := filter["status"]; ok && status == StatusActive {
		tx.Where("status = ? OR status = ? OR status", StatusCreate, StatusCreate, StatusCreate)
		delete(filter, "status")
	}

	var orders []Order
	if err := tx.
		Where(filter).
		Find(&orders).Error; err != nil {
		shared.LogError("error finding orders", LogDBRepository, "Find", err, filter)
		return nil, err
	}

	return orders, nil
}

// UpdateOrderItem method for update an order item in database
func (r *DBRepository) UpdateOrderItem(item *OrderItem) (*OrderItem, error) {
	var orderItem OrderItem
	if err := r.db.First(&orderItem, item.ID).Error; err != nil {
		shared.LogError("error getting order item", LogDBRepository, "UpdateOrderItem", err, *item)
		return nil, err
	}

	if err := r.db.Model(&orderItem).Updates(item).Error; err != nil {
		shared.LogError("error updating order item", LogDBRepository, "UpdateOrderItem", err, *item)
		return nil, err
	}

	return &orderItem, nil
}

// GetOrderItem method for get an order item from database
func (r *DBRepository) GetOrderItem(orderItemID string) (*OrderItem, error) {
	var orderItem OrderItem
	if err := r.db.First(&orderItem, orderItemID).Error; err != nil {
		shared.LogError("error getting order item", LogDBRepository, "GetOrderItem", err, orderItemID)
		return nil, err
	}

	return &orderItem, nil
}

// OrderType methods

// CreateOrderType method for create a new order type in database
func (r *DBRepository) CreateOrderType(orderType *OrderType) (*OrderType, error) {
	if err := r.db.Save(orderType).Error; err != nil {
		shared.LogError("error creating order type", LogDBRepository, "CreateOrderType", err, *orderType)
		return nil, fmt.Errorf(ErrorOrderTypeCreation)
	}

	return orderType, nil
}

// FindOrderType method for find order types in database
func (r *DBRepository) FindOrderType(filter map[string]any) ([]OrderType, error) {
	var orderTypes []OrderType
	if err := r.db.Find(&orderTypes, filter).Error; err != nil {
		shared.LogError("error finding order types", LogDBRepository, "FindOrderType", err, filter)
		return nil, fmt.Errorf(ErrorOrderTypeFinding)
	}

	return orderTypes, nil
}

// GetOrderType method for get an order type from database
func (r *DBRepository) GetOrderType(orderTypeID string) (*OrderType, error) {
	var orderType OrderType
	if err := r.db.First(&orderType, orderTypeID).Error; err != nil {
		shared.LogError("error getting order type", LogDBRepository, "GetOrderType", err, orderTypeID)
		return nil, fmt.Errorf(ErrorOrderTypeGetting)
	}

	return &orderType, nil
}

// UpdateOrderType method for update an order type in database
func (r *DBRepository) UpdateOrderType(orderTypeID string, orderType *OrderType) (*OrderType, error) {
	var orderTypeDB OrderType
	if err := r.db.First(&orderTypeDB, orderTypeID).Error; err != nil {
		shared.LogError("error getting order type", LogDBRepository, "UpdateOrderType", err, orderTypeID)
		return nil, fmt.Errorf(ErrorOrderTypeGetting)
	}

	if err := r.db.Model(&orderTypeDB).Updates(orderType).Error; err != nil {
		shared.LogError("error updating order type", LogDBRepository, "UpdateOrderType", err, *orderType)
		return nil, fmt.Errorf(ErrorOrderTypeUpdating)
	}

	return &orderTypeDB, nil
}

// DeleteOrderType method for delete an order type in database
func (r *DBRepository) DeleteOrderType(orderTypeID string) error {
	if err := r.db.Delete(&OrderType{}, orderTypeID).Error; err != nil {
		shared.LogError("error deleting order type", LogDBRepository, "DeleteOrderType", err, orderTypeID)
		return fmt.Errorf(ErrorOrderTypeDeleting)
	}

	return nil
}

// Attendee methods

func (r *DBRepository) CreateAttendee(attendee *Attendee) (*Attendee, error) {
	if err := r.db.Save(attendee).Error; err != nil {
		shared.LogError("error creating attendee", LogDBRepository, "CreateAttendee", err, *attendee)
		return nil, err
	}

	return attendee, nil
}
