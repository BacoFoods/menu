package order

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/BacoFoods/menu/pkg/shared"

	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/client"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/tables"
	"gorm.io/gorm"
)

const (
	ErrorBadRequest                        = "error bad request"
	ErrorBadRequestOrderID                 = "error bad request wrong order id"
	ErrorBadRequestOrderItemID             = "error bad request wrong order item id"
	ErrorBadRequestProductID               = "error bad request wrong product id"
	ErrorBadRequestTableID                 = "error bad request wrong table id"
	ErrorBadRequestStoreID                 = "error bad request wrong store id"
	ErrorBadRequestOrderSeats              = "error bad request wrong order seats can't be less than 0"
	ErrorOrderCreation                     = "error creating order"
	ErrorOrderGetting                      = "error getting order"
	ErrorOrderFind                         = "error finding orders"
	ErrorOrderAddProductsForbiddenByStatus = "error adding products to order forbidden by status"
	ErrorOrderUpdate                       = "error updating order"
	ErrorOrderUpdateStatus                 = "error updating order status"
	ErrorOrderUpdateInvalidStatus          = "error updating order invalid status"
	ErrorOrderProductGetting               = "error getting order product"
	ErrorOrderProductNotFound              = "error order product with id %v not found; "
	ErrorOrderProductsNotFound             = "error order products not found"
	ErrorOrderModifierNotFound             = "error order modifier with id %v not found; "
	ErrorOrderUpdatingComments             = "error updating order comments"
	ErrorOrderUpdatingClientName           = "error updating order client name"
	ErrorOrderInvoiceCreation              = "error creating order invoice"
	ErrorOrderInvoiceCalculation           = "error calculating invoice"

	ErrorOrderItemUpdate       = "error updating order item"
	ErrorOrderItemGetting      = "error getting order item"
	ErrorOrderItemUpdateCourse = "error updating order item course"

	ErrorOrderTypeCreation = "error creating order type"
	ErrorOrderTypeFinding  = "error finding order type"
	ErrorOrderTypeGetting  = "error getting order type"
	ErrorOrderTypeUpdating = "error updating order type"
	ErrorOrderTypeDeleting = "error deleting order type"

	TaxPercentage = 0.08

	OrderStepCreated OrderStep = "created"

	OrderActionCreated OrderAction = "fue atendido por"

	LogDomain = "pkg/order/domain"

	OrderStatusCreated = "created"
	OrderStatusPaying  = "paying"
	OrderStatusClosed  = "closed"
)

type OrderStep string

type OrderAction string

type Repository interface {
	// Order
	Create(order *Order) (*Order, error)
	Get(orderID string) (*Order, error)
	Find(filter map[string]any) ([]Order, error)
	Update(order *Order) (*Order, error)
	FindByShift(shiftID uint) ([]Order, error)
	AddProducts(order *Order, newItems []OrderItem) (*Order, error)

	// OrderItem
	UpdateOrderItem(orderItem *OrderItem) (*OrderItem, error)
	GetOrderItem(orderItemID string) (*OrderItem, error)

	// OrderType
	CreateOrderType(*OrderType) (*OrderType, error)
	FindOrderType(filter map[string]any) ([]OrderType, error)
	GetOrderType(orderTypeID string) (*OrderType, error)
	UpdateOrderType(orderTypeID string, orderType *OrderType) (*OrderType, error)
	DeleteOrderType(orderTypeID string) error

	// Attendee
	CreateAttendee(attendee *Attendee) (*Attendee, error)
}

type Order struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	Statuses      []OrderStatus     `json:"status" gorm:"foreignKey:OrderID" swaggerignore:"true"`
	CurrentStatus string            `json:"current_status"`
	OrderType     string            `json:"order_type"`
	ClientName    string            `json:"client_name"`
	BrandID       *uint             `json:"brand_id" binding:"required"`
	Brand         *brand.Brand      `json:"brand,omitempty" swaggerignore:"true"`
	StoreID       *uint             `json:"store_id" binding:"required"`
	Store         *store.Store      `json:"store,omitempty" swaggerignore:"true"`
	ChannelID     *uint             `json:"channel_id" binding:"required"`
	TableID       *uint             `json:"table_id"`
	Table         *tables.Table     `json:"table,omitempty" swaggerignore:"true"`
	TypeID        *uint             `json:"type_id"`
	Type          *OrderType        `json:"type"`
	Comments      string            `json:"comments"`
	Items         []OrderItem       `json:"items"  gorm:"foreignKey:OrderID"`
	CookingTime   int               `json:"cooking_time"`
	Seats         int               `json:"seats"`
	ExternalCode  string            `json:"external_code"`
	Invoices      []invoice.Invoice `json:"invoices"  gorm:"foreignKey:OrderID" swaggerignore:"true"`
	Attendees     []Attendee        `json:"attendees" gorm:"foreignKey:OrderID"`
	ShiftID       *uint             `json:"shift_id"`
	CreatedAt     *time.Time        `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt     *time.Time        `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt     *gorm.DeletedAt   `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (o *Order) GetProductIDs() []string {
	ids := make([]string, len(o.Items))
	for i, item := range o.Items {
		ids[i] = fmt.Sprintf("%d", *item.ProductID)
	}
	return ids
}

func (o *Order) GetModifierIDs() []string {
	ids := make([]string, 0)
	for _, item := range o.Items {
		for _, modifier := range item.Modifiers {
			ids = append(ids, fmt.Sprintf("%d", *modifier.ProductID))
		}
	}
	return ids
}

func (o *Order) SetItems(products []product.Product, modifiers []product.Product) {
	productsMap := make(map[string]product.Product)
	for _, p := range products {
		productsMap[fmt.Sprintf("%d", p.ID)] = p
	}

	modifiersMap := make(map[string]product.Product)
	for _, m := range modifiers {
		modifiersMap[fmt.Sprintf("%d", m.ID)] = m
	}

	items := make([]OrderItem, 0)
	for _, item := range o.Items {
		if p, ok := productsMap[fmt.Sprintf("%d", *item.ProductID)]; ok {
			item.Name = p.Name
			item.Description = p.Description
			item.Image = p.Image
			item.SKU = p.SKU
			item.Price = p.Price
			item.Unit = p.Unit
			// TODO: add default tax value if Tax is nil
			if p.Tax != nil {
				item.Tax = p.Tax.Name
				item.TaxPercentage = p.Tax.Percentage
			}

			item.SetHash()

			modifierList := make([]OrderModifier, 0)
			for _, modifier := range item.Modifiers {
				if m, ok := modifiersMap[fmt.Sprintf("%d", *modifier.ProductID)]; ok {
					modifier.Name = m.Name
					modifier.Description = m.Description
					modifier.Image = m.Image
					modifier.SKU = m.SKU
					modifier.Price = m.Price
					modifier.Unit = m.Unit
					modifier.ProductID = &m.ID
					modifier.OrderID = o.ID
					modifierList = append(modifierList, modifier)
				}
			}
			item.Modifiers = modifierList
			items = append(items, item)
		}
	}
	o.Items = items
}

func (o *Order) AddProduct(orderItem OrderItem) {
	o.Items = append(o.Items, orderItem)
}

func (o *Order) RemoveProduct(product *product.Product) {
	for i, item := range o.Items {
		if *item.ProductID == product.ID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			return
		}
	}
}

func (o *Order) ToInvoice(tip *tipData) {
	// Remove invoices
	o.Invoices = nil
	subtotal := 0.0
	newInvoice := invoice.Invoice{
		OrderID:   &o.ID,
		BrandID:   o.BrandID,
		StoreID:   o.StoreID,
		ChannelID: o.ChannelID,
		TableID:   o.TableID,
		ShiftID:   o.ShiftID,
		Items:     make([]invoice.Item, 0),
		Client:    client.DefaultClient(),
		BaseTax:   0,
	}

	// Adding items to invoice
	for _, item := range o.Items {
		tax := "ico" // Default tax
		if item.Tax != "" {
			tax = item.Tax
		}

		taxPerc := 0.08
		if item.TaxPercentage != 0 {
			taxPerc = item.TaxPercentage
		}

		productBaseTax := math.Floor(item.Price / (1 + taxPerc))
		newInvoice.BaseTax += productBaseTax
		newInvoice.Taxes += item.Price - productBaseTax

		newInvoice.Items = append(newInvoice.Items, invoice.Item{
			ProductID:     item.ProductID,
			Name:          item.Name,
			Description:   item.Description,
			SKU:           item.SKU,
			Price:         item.Price,
			Comments:      item.Comments,
			Hash:          item.Hash,
			Tax:           tax,
			TaxPercentage: taxPerc,
		})

		// Adding item price to subtotal
		subtotal += item.Price

		for _, modifier := range item.Modifiers {
			// Adding modifier price to subtotal
			subtotal += modifier.Price

			modifierTax := "ico" // Default tax
			if modifier.Tax != "" {
				modifierTax = modifier.Tax
			}

			modifierTaxPerc := 0.08
			if modifier.TaxPercentage != 0 {
				modifierTaxPerc = modifier.TaxPercentage
			}

			modifierBaseTax := math.Floor(modifier.Price / (1 + modifierTaxPerc))
			newInvoice.BaseTax += modifierBaseTax
			newInvoice.Taxes += modifier.Price - modifierBaseTax

			newInvoice.Items = append(newInvoice.Items, invoice.Item{
				ProductID:     modifier.ProductID,
				Name:          modifier.Name,
				Description:   modifier.Description,
				SKU:           modifier.SKU,
				Price:         modifier.Price,
				Comments:      modifier.Comments,
				Tax:           modifierTax,
				TaxPercentage: modifierTaxPerc,
			})
		}
	}

	// Setting subtotals
	newInvoice.SubTotal = subtotal

	// Setting taxes
	newInvoice.CalculateTaxDetails()

	if tip != nil {
		tipType, tipValue := tip.GetValueAndType()
		newInvoice.Tip = fmt.Sprintf("%s - %f", tipType, tipValue)
		if tipType == "percentage" {
			tipAmount := math.Floor(newInvoice.BaseTax * tipValue)
			newInvoice.TipAmount = tipAmount
		} else {
			newInvoice.TipAmount = tipValue
		}
	}

	newInvoice.Total = newInvoice.SubTotal + newInvoice.TipAmount

	// Setting invoice
	o.Invoices = []invoice.Invoice{newInvoice}
}

func (o *Order) UpdateStatus(status string) error {
	if status == o.CurrentStatus {
		shared.LogWarn("order already has this status", LogDomain, "UpdateStatus", nil, o.ID, o.CurrentStatus, status)
		return nil
	}

	if status == OrderStatusCreated && o.CurrentStatus == "" {
		o.CurrentStatus = status
		o.Statuses = append(o.Statuses, OrderStatus{
			Code:    OrderStatusCreated,
			OrderID: &o.ID,
		})
		return nil
	}

	if status == OrderStatusPaying && o.CurrentStatus == OrderStatusCreated {
		o.CurrentStatus = status
		o.Statuses = append(o.Statuses, OrderStatus{
			Code:    OrderStatusPaying,
			OrderID: &o.ID,
		})
		return nil
	}

	if status == OrderStatusClosed && o.CurrentStatus == OrderStatusPaying {
		o.CurrentStatus = status
		o.Statuses = append(o.Statuses, OrderStatus{
			Code:    OrderStatusClosed,
			OrderID: &o.ID,
		})
		return nil
	}

	err := fmt.Errorf(ErrorOrderUpdateStatus)
	shared.LogError(ErrorOrderUpdateStatus, LogDomain, "UpdateStatus", nil, o.ID, o.CurrentStatus, status)
	return err
}

func (o *Order) UpdateNextStatus() {
	currentStatus := o.Statuses[len(o.Statuses)-1]
	o.Statuses = append(o.Statuses, currentStatus.Next())
	o.CurrentStatus = currentStatus.Next().Code
}

func (o *Order) UpdatePrevStatus() {
	currentStatus := o.Statuses[len(o.Statuses)-1]
	o.Statuses = append(o.Statuses, currentStatus.Prev())
	o.CurrentStatus = currentStatus.Prev().Code
}

type OrderItem struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	OrderID         *uint           `json:"order_id"`
	ProductID       *uint           `json:"product_id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Image           string          `json:"image"`
	SKU             string          `json:"sku"`
	Price           float64         `json:"price" gorm:"precision:18;scale:2"`
	Unit            string          `json:"unit"`
	Discount        float64         `json:"discount" gorm:"precision:18;scale:2"`
	DiscountReason  string          `json:"discount_reason,omitempty"`
	Surcharge       float64         `json:"surcharge" gorm:"precision:18;scale:2"`
	SurchargeReason string          `json:"surcharge_reason,omitempty"`
	Comments        string          `json:"comments"`
	Course          string          `json:"course"`
	Hash            string          `json:"hash"`
	Modifiers       []OrderModifier `json:"modifiers"  gorm:"foreignKey:OrderItemID"`
	Tax             string          `json:"tax"`
	TaxPercentage   float64         `json:"tax_percentage"`
	CreatedAt       *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (oi *OrderItem) SetHash() {
	orderItemString := fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v", oi.ID, oi.OrderID, oi.Name, oi.Description, oi.Image, oi.Price, oi.Unit, oi.Discount, oi.DiscountReason, oi.Surcharge, oi.SurchargeReason, oi.Comments, oi.Course)
	for _, modifier := range oi.Modifiers {
		orderItemString += fmt.Sprintf("%v%v%v%v%v%v%v%v", modifier.ID, modifier.OrderItemID, modifier.Name, modifier.Description, modifier.Image, modifier.Price, modifier.Unit, modifier.Comments)
	}
	oi.Hash = fmt.Sprintf("%x", orderItemString)
}

func (oi *OrderItem) AddModifiers(modifier []OrderModifier) {
	oi.Modifiers = append(oi.Modifiers, modifier...)
}

func (oi *OrderItem) RemoveModifiers(modifiers []OrderModifier) {
	for _, modifier := range modifiers {
		for i, m := range oi.Modifiers {
			if m.ID == modifier.ID {
				oi.Modifiers = append(oi.Modifiers[:i], oi.Modifiers[i+1:]...)
				break
			}
		}
	}
}

type OrderModifier struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	OrderItemID   *uint           `json:"order_item_id"`
	OrderID       uint            `json:"order_id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Image         string          `json:"image"`
	Category      string          `json:"category"`
	ProductID     *uint           `json:"product_id"`
	SKU           string          `json:"sku"`
	Price         float64         `json:"price"  gorm:"precision:18;scale:2"`
	Unit          string          `json:"unit"`
	Tax           string          `json:"tax"`
	TaxPercentage float64         `json:"tax_percentage"`
	Comments      string          `json:"comments"`
	CreatedAt     *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt     *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt     *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type OrderType struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	ChannelID   *uint           `json:"channel_id"`
	StoreID     *uint           `json:"store_id"`
	BrandID     *uint           `json:"brand_id"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Attendee struct {
	ID        uint            `json:"id" swaggerignore:"true" gorm:"primaryKey"`
	AccountID uint            `json:"account_id"`
	OrderID   uint            `json:"order_id" swaggerignore:"true"`
	Name      string          `json:"name"`
	Role      string          `json:"role"`
	Action    OrderAction     `json:"order_action"`
	OrderStep OrderStep       `json:"order_step"`
	CreatedAt *time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`
}

type OrderStatus struct {
	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
	Code      string         `json:"code"`
	OrderID   *uint          `json:"order_id,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (os *OrderStatus) Next() OrderStatus {
	switch os.Code {
	case OrderStatusCreated:
		return OrderStatus{
			Code: OrderStatusPaying,
		}
	case OrderStatusPaying:
		return OrderStatus{
			Code: OrderStatusClosed,
		}
	case OrderStatusClosed:
		return OrderStatus{
			Code: OrderStatusClosed,
		}
	default:
		return OrderStatus{
			Code: OrderStatusCreated,
		}
	}
}

func (os *OrderStatus) Prev() OrderStatus {
	switch os.Code {
	case OrderStatusCreated:
		return OrderStatus{
			Code: OrderStatusCreated,
		}
	case OrderStatusPaying:
		return OrderStatus{
			Code: OrderStatusCreated,
		}
	case OrderStatusClosed:
		return OrderStatus{
			Code: OrderStatusPaying,
		}
	default:
		return OrderStatus{
			Code: OrderStatusCreated,
		}
	}
}

type tipData struct {
	Percentage string `json:"percentage"`
	Amount     string `json:"value"`
}

func (t tipData) GetValueAndType() (string, float64) {
	if t.Percentage != "" {
		p, err := strconv.ParseFloat(t.Percentage, 64)
		if err == nil && p > 0 {
			return "percentage", p / 100
		}
	}

	if t.Amount != "" {
		v, err := strconv.ParseFloat(t.Amount, 64)
		if err == nil && v > 0 {
			return "value", v
		}
	}

	return "", 0
}
