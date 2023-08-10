package order

import (
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/tables"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest    = "error bad request"
	ErrorOrderCreation = "error creating order"
)

type Order struct {
	ID            uint             `json:"id" gorm:"primaryKey"`
	Statuses      []OrderStatus    `json:"status" gorm:"foreignKey:OrderID" swaggerignore:"true"`
	CurrentStatus string           `json:"current_status"`
	OrderType     string           `json:"order_type"`
	BrandID       *uint            `json:"brand_id" binding:"required"`
	StoreID       *uint            `json:"store_id" binding:"required"`
	ChannelID     *uint            `json:"channel_id" binding:"required"`
	TableID       *uint            `json:"table_id"`
	Table         *tables.Table    `json:"table"`
	Comments      string           `json:"comments"`
	Items         []OrderItem      `json:"items"  gorm:"foreignKey:OrderID"`
	Discounts     []OrderDiscount  `json:"discounts"  gorm:"foreignKey:OrderID"`
	Surcharges    []OrderSurcharge `json:"surcharges"  gorm:"foreignKey:OrderID"`
	CookingTime   int              `json:"cooking_time"`
	Seats         int              `json:"seats"`
	ExternalCode  string           `json:"external_code"`
	InvoiceID     *uint            `json:"invoice_id"`
	Invoice       *invoice.Invoice `json:"invoice"`
	CreatedAt     *time.Time       `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt     *time.Time       `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt     *gorm.DeletedAt  `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type OrderItem struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	OrderID         *uint           `json:"order_id"`
	ProductID       *uint           `json:"product_id"`
	Product         product.Product `json:"product"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Image           string          `json:"image"`
	SKU             string          `json:"sku"`
	Price           float64         `json:"price" gorm:"precision:18;scale:2"`
	Unit            string          `json:"unit"`
	Quantity        int             `json:"quantity"`
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

type OrderModifier struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	OrderItemID *uint           `json:"order_item_id"`
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

type OrderStatus struct {
	ID          uint            `json:"id"`
	OrderID     *uint           `json:"order_id"`
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

type OrderType struct {
	ID          uint            `json:"id"`
	OrderID     *uint           `json:"order_id"`
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

type OrderDiscount struct {
	ID          uint           `json:"id"`
	OrderID     *uint          `json:"order_id"`
	Name        string         `json:"name,omitempty"`
	Type        string         `json:"type"`
	Percentage  float32        `json:"percentage,omitempty" gorm:"precision:18;scale:2"`
	Value       float32        `json:"value,omitempty" gorm:"precision:18;scale:2"`
	Description string         `json:"description,omitempty"`
	Terms       string         `json:"terms,omitempty"`
	ChannelID   *uint          `json:"channel_id,omitempty"`
	StoreID     *uint          `json:"store_id,omitempty"`
	BrandID     *uint          `json:"brand_id,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type OrderSurcharge struct {
	ID          uint            `json:"id"`
	OrderID     *uint           `json:"order_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Percentage  float64         `json:"percentage" gorm:"precision:18;scale:2"`
	Amount      float64         `json:"amount" gorm:"precision:18;scale:2"`
	Active      bool            `json:"active"`
	ChannelID   *uint           `json:"channel_id,omitempty"`
	StoreID     *uint           `json:"store_id,omitempty"`
	BrandID     *uint           `json:"brand_id,omitempty"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(order *Order) (*Order, error)
	Get(orderID string) (*Order, error)
}
