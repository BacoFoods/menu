package order

import (
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/surcharge"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	Status      []Status              `json:"status"`
	TypeID      *uint                 `json:"type_id"`
	Type        Type                  `json:"type"`
	BrandID     *uint                 `json:"brand_id"`
	StoreID     *uint                 `json:"store_id"`
	ChannelID   *uint                 `json:"channel_id"`
	TableID     *uint                 `json:"table_id"`
	Comments    string                `json:"comments"`
	Detail      []OrderDetail         `json:"items"`
	Discounts   []discount.Discount   `json:"discounts"`
	Surcharges  []surcharge.Surcharge `json:"surcharges"`
	CookingTime int                   `json:"cooking_time"`
	Eaters      int                   `json:"eaters"`
	InvoiceID   *uint                 `json:"invoice_id"`
	CreatedAt   *time.Time            `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time            `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt       `json:"deleted_at,omitempty" swaggerignore:"true"`
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
	Comments        string          `json:"comments"`
	Course          Course          `json:"course"`
	CreatedAt       *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Course string

type Status string

const (
	InProgress Status = "IN_PROGRESS"
	Ready      Status = "READY"
	Delivered  Status = "DELIVERED"
	Canceled   Status = "CANCELED"

	Entrance Course = "ENTRANCE"
	Shared   Course = "SHARED"
	Main     Course = "MAIN"
)

type Type string

const (
	TakeAway     Type = "TAKE_AWAY"
	QuickService Type = "QUICK_SERVICE"
	FullService  Type = "FULL_SERVICE"
	Delivery     Type = "DELIVERY"
)

type Repository interface {
	Find(filters map[string]string) ([]Order, error)
}
