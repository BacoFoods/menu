package order

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/status"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/tables"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest            = "error bad request"
	ErrorBadRequestOrderID     = "error bad request wrong order id"
	ErrorBadRequestOrderItemID = "error bad request wrong order item id"
	ErrorBadRequestProductID   = "error bad request wrong product id"
	ErrorBadRequestTableID     = "error bad request wrong table id"
	ErrorBadRequestStoreID     = "error bad request wrong store id"
	ErrorBadRequestOrderSeats  = "error bad request wrong order seats can't be less than 0"
	ErrorOrderCreation         = "error creating order"
	ErrorOrderGetting          = "error getting order"
	ErrorOrderUpdate           = "error updating order"
	ErrorOrderProductGetting   = "error getting order product"
	ErrorOrderProductNotFound  = "error order product with id %v not found; "
	ErrorOrderModifierNotFound = "error order modifier with id %v not found; "

	ErrorOrderItemUpdate  = "error updating order item"
	ErrorOrderItemGetting = "error getting order item"

	ErrorOrderTypeCreation = "error creating order type"
	ErrorOrderTypeFinding  = "error finding order type"
	ErrorOrderTypeGetting  = "error getting order type"
	ErrorOrderTypeUpdating = "error updating order type"
	ErrorOrderTypeDeleting = "error deleting order type"

	TaxPercentage = 0.08
)

type Repository interface {
	Create(order *Order) (*Order, error)
	Get(orderID string) (*Order, error)
	Find(filter map[string]any) ([]Order, error)
	Update(order *Order) (*Order, error)
	UpdateOrderItem(orderItem *OrderItem) (*OrderItem, error)
	GetOrderItem(orderItemID string) (*OrderItem, error)

	CreateOrderType(*OrderType) (*OrderType, error)
	FindOrderType(filter map[string]any) ([]OrderType, error)
	GetOrderType(orderTypeID string) (*OrderType, error)
	UpdateOrderType(orderTypeID string, orderType *OrderType) (*OrderType, error)
	DeleteOrderType(orderTypeID string) error
}

type Order struct {
	ID            uint             `json:"id" gorm:"primaryKey"`
	Statuses      []status.Status  `json:"status" gorm:"many2many:order_statuses" gorm:"foreignKey:OrderID" swaggerignore:"true"`
	CurrentStatus string           `json:"current_status"`
	OrderType     string           `json:"order_type"`
	BrandID       *uint            `json:"brand_id" binding:"required"`
	Brand         *brand.Brand     `json:"brand,omitempty"`
	StoreID       *uint            `json:"store_id" binding:"required"`
	Store         *store.Store     `json:"store,omitempty"`
	ChannelID     *uint            `json:"channel_id" binding:"required"`
	TableID       *uint            `json:"table_id"`
	Table         *tables.Table    `json:"table"`
	TypeID        *uint            `json:"type_id"`
	Type          *OrderType       `json:"type"`
	Comments      string           `json:"comments"`
	Items         []OrderItem      `json:"items"  gorm:"foreignKey:OrderID"`
	CookingTime   int              `json:"cooking_time"`
	Seats         int              `json:"seats"`
	ExternalCode  string           `json:"external_code"`
	InvoiceID     *uint            `json:"invoice_id"`
	Invoice       *invoice.Invoice `json:"invoice"`
	CreatedAt     *time.Time       `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt     *time.Time       `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt     *gorm.DeletedAt  `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (o *Order) GetProductIDs() []string {
	ids := make([]string, len(o.Items))
	for i, item := range o.Items {
		ids[i] = fmt.Sprintf("%d", *item.ProductID)
	}
	return ids
}

func (o *Order) SetItems(products []product.Product) {
	productsMap := make(map[string]product.Product)
	for _, p := range products {
		productsMap[fmt.Sprintf("%d", p.ID)] = p
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

func (o *Order) ToInvoice() {
	subtotal := 0.0

	newInvoice := invoice.Invoice{
		OrderID:   &o.ID,
		BrandID:   o.BrandID,
		StoreID:   o.StoreID,
		ChannelID: o.ChannelID,
		TableID:   o.TableID,
		Table:     o.Table,
		Items:     make([]invoice.Item, 0),
	}

	// Adding items to invoice
	for _, item := range o.Items {
		newInvoice.Items = append(newInvoice.Items, invoice.Item{
			ProductID:   item.ProductID,
			Name:        item.Name,
			Description: item.Description,
			SKU:         item.SKU,
			Price:       item.Price,
			Comments:    item.Comments,
		})

		// Adding item price to subtotal
		subtotal += item.Price

		for _, modifier := range item.Modifiers {
			// Only modifiers with price are added to invoice
			if modifier.Price == 0 {
				continue
			}

			// Adding modifier price to subtotal
			subtotal += modifier.Price

			newInvoice.Items = append(newInvoice.Items, invoice.Item{
				ProductID:   modifier.ProductID,
				Name:        modifier.Name,
				Description: modifier.Description,
				SKU:         modifier.SKU,
				Price:       modifier.Price,
				Comments:    modifier.Comments,
			})
		}
	}

	// Setting subtotals
	newInvoice.SubTotal = subtotal

	// Setting taxes
	tax := subtotal * TaxPercentage
	baseTax := subtotal - tax
	newInvoice.BaseTax = baseTax
	newInvoice.Total = subtotal + tax

	// Setting invoice
	o.Invoice = &newInvoice
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
	Modifiers       []OrderModifier `json:"modifiers"  gorm:"foreignKey:OrderItemID"`
	CreatedAt       *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
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
	ID          uint            `json:"id" gorm:"primaryKey"`
	OrderItemID *uint           `json:"order_item_id"`
	OrderID     uint            `json:"order_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	Category    string          `json:"category"`
	ProductID   *uint           `json:"product_id"`
	SKU         string          `json:"sku"`
	Price       float64         `json:"price"  gorm:"precision:18;scale:2"`
	Unit        string          `json:"unit"`
	Comments    string          `json:"comments"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
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
