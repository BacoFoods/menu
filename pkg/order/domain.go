package order

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	Status      []Status        `json:"status"`
	TypeID      *uint           `json:"type_id"`
	Type        Type            `json:"type"`
	BrandID     *uint           `json:"brand_id"`
	StoreID     *uint           `json:"store_id"`
	ChannelID   *uint           `json:"channel_id"`
	TableID     *uint           `json:"table_id"`
	Observation string          `json:"observation"`
	Detail      []OrderDetail   `json:"items"`
	Discounts   []Discount      `json:"discounts"`
	Surcharges  []Surcharge     `json:"surcharges"`
	CookingTime int             `json:"cooking_time"`
	Eaters      int             `json:"eaters"`
	InvoiceID   *uint           `json:"invoice_id"`
	Invoice     Invoice         `json:"invoice"`
	ClientID    *uint           `json:"client_id"` // optional
	Client      Client          `json:"client"`    // optional
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type OrderDetail struct {
	OrderID         *uint           `json:"order_id"`
	ProductID       *uint           `json:"product_id"`
	Quantity        int             `json:"quantity"`
	Price           float64         `json:"price"`
	Discount        float64         `json:"discount"`
	DiscountReason  string          `json:"discount_reason"`
	Surcharge       float64         `json:"surcharge"`
	SurchargeReason string          `json:"surcharge_reason"`
	Observation     string          `json:"observation"`
	CreatedAt       *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Status string

const (
	InProgress Status = "IN_PROGRESS"
	Ready      Status = "READY"
	Delivered  Status = "DELIVERED"
	Canceled   Status = "CANCELED"
)

type Type string

const (
	TakeAway     Type = "TAKE_AWAY"
	QuickService Type = "QUICK_SERVICE"
	FullService  Type = "FULL_SERVICE"
	Delivery     Type = "DELIVERY"
)
