package order

import (
	"fmt"
	"strings"

	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	LogDBRepository string = "pkg/order/db_repository"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

// Order methods

// Create method for create a new order in database
func (r *DBRepository) Create(order *Order, ch *channel.Channel) (*Order, error) {
	isNew := order.ID == 0
	if err := r.db.Save(order).Error; err != nil {
		shared.LogError("error creating order", LogDBRepository, "Create", err, *order)
		return nil, err
	}

	if isNew {
		ch := strings.ToUpper(ch.ShortName)
		order.Code = fmt.Sprintf("%s%d", ch, order.ID)
	}

	if err := r.db.Save(order).Error; err != nil {
		shared.LogError("error updating order code", LogDBRepository, "Create", err, *order)
		return nil, err
	}

	return order, nil
}

// Get method for get an order from database
func (r *DBRepository) Get(orderID string) (*Order, error) {
	if strings.TrimSpace(orderID) == "" {
		err := fmt.Errorf(ErrorOrderIDEmpty)
		shared.LogWarn("error getting order", LogDBRepository, "Get", err)
		return nil, err
	}

	var order Order
	if err := r.db.
		Preload(clause.Associations).
		Preload("Invoices.Documents").
		Preload("Items.Modifiers").
		Preload("Table.Zone").
		First(&order, orderID).Error; err != nil {
		shared.LogError("error getting order", LogDBRepository, "Get", err, orderID)
		return nil, err
	}

	return &order, nil
}

// AddProducts method for add products to an order in database
func (c *DBRepository) AddProducts(order *Order, newItems []OrderItem) (*Order, error) {
	if err := c.db.Model(order).
		Association("Items").
		Append(newItems); err != nil {
		return nil, err
	}

	if err := c.db.Save(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

// Update method for update an order in database
func (r *DBRepository) Update(order *Order) (*Order, error) {
	if err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(order).Error; err != nil {
		shared.LogError("error updating order", LogDBRepository, "Update", err, *order)
		return nil, err
	}

	return order, nil
}

// UpdateTable method for update an order table in database
func (r *DBRepository) UpdateTable(order *Order, newTableID uint) (*Order, error) {
	return order, r.db.Model(order).Select("table_id").Where("id = ?", order.ID).Update("table_id", newTableID).Error
}

// Find method for find orders in database
func (r *DBRepository) Find(filter map[string]any) ([]Order, error) {
	tx := r.db.
		Preload(clause.Associations).
		Preload("Items.Modifiers").
		Preload("Table.Zone").
		Preload("Invoices.Payments")

	if days, ok := filter["days"]; ok {
		tx.Preload(clause.Associations).
			Preload("Items.Modifiers").
			Where(fmt.Sprintf("created_at >= NOW() - INTERVAL '%s' DAY", days))

		shared.LogWarn("filtering by days", LogDBRepository, "Find", nil, filter)
		delete(filter, "days")
	}

	if statuses, ok := filter["current_status"]; ok {
		tx.Where("current_status IN ?", statuses)
		delete(filter, "current_status")
	}

	var orders []Order
	err := tx.
		Preload(clause.Associations).
		Preload("Items.Modifiers").
		Where(filter).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
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

// FindByShift method for find orders by shift in database
func (r *DBRepository) FindByShift(shiftID uint) ([]Order, error) {
	var orders []Order
	if err := r.db.Preload(clause.Associations).Find(&orders, "shift_id = ?", shiftID).Error; err != nil {
		shared.LogError("error finding orders", LogDBRepository, "FindByShift", err, shiftID)
		return nil, err
	}

	return orders, nil
}

// GetLastDayOrders method for get day's orders from database
func (r *DBRepository) GetLastDayOrders(storeID string) ([]Order, error) {
	var orders []Order
	if err := r.db.Preload(clause.Associations).
		Preload("Invoices.Payments").
		Where("store_id = ? AND created_at >= NOW() - INTERVAL '1' DAY", storeID).
		Find(&orders).Error; err != nil {
		shared.LogError("error getting orders", LogDBRepository, "GetLastDayOrders", err, storeID)
		return nil, err
	}

	return orders, nil
}

// GetLastDayOrdersByStatus method for get day's orders from database
func (r *DBRepository) GetLastDayOrdersByStatus(storeID string, status string) ([]Order, error) {
	var orders []Order
	if err := r.db.Preload(clause.Associations).
		Preload("Invoices.Payments").
		Where("store_id = ? AND current_status = ? AND created_at >= NOW() - INTERVAL '1' DAY", storeID, status).
		Find(&orders).Error; err != nil {
		shared.LogError("error getting orders", LogDBRepository, "GetLastDayOrders", err, storeID)
		return nil, err
	}

	return orders, nil
}

// Delete method for delete an order in database
func (r *DBRepository) Delete(orderID string) error {
	if err := r.db.Delete(&Order{}, orderID).Error; err != nil {
		shared.LogError("error deleting order", LogDBRepository, "Delete", err, orderID)
		return fmt.Errorf(ErrorOrderDeleting)
	}

	return nil
}

func (r *DBRepository) FindByIdempotencyKey(idempotencyKey string, storeID *uint) (*Order, error) {
	var order Order

	if err := r.db.Preload(clause.Associations).
		Where("idempotency_key = ?", idempotencyKey).
		Where("store_id = ?", storeID).
		First(&order).Error; err != nil {
		shared.LogError("error getting order", LogDBRepository, "FindByIdempotencyKey", err, idempotencyKey)
		return nil, err
	}

	return &order, nil
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

var _ Repository = (*DBRepository)(nil)
