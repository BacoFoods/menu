package order

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/order/db_repository"

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
	var orders []Order
	if err := r.db.
		Preload(clause.Associations).
		Preload("Items.Modifiers").
		Where(filter).
		Find(&orders).Error; err != nil {
		shared.LogError("error finding orders", LogDBRepository, "Find", err, filter)
		return nil, err
	}

	return orders, nil
}
