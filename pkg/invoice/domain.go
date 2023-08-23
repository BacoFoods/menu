package invoice

import (
	"github.com/BacoFoods/menu/pkg/payment"
	"github.com/BacoFoods/menu/pkg/tables"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest      = "error bad request"
	ErrorInvoiceCreation = "error creating invoice"
)

type Repository interface {
	Create(invoice *Invoice) (*Invoice, error)
}

type Invoice struct {
	ID              uint             `json:"id"`
	OrderID         *uint            `json:"order_id"`
	BrandID         *uint            `json:"brand_id" binding:"required"`
	StoreID         *uint            `json:"store_id" binding:"required"`
	ChannelID       *uint            `json:"channel_id" binding:"required"`
	TableID         *uint            `json:"table_id"`
	Table           *tables.Table    `json:"table"`
	Items           []Item           `json:"items"  gorm:"foreignKey:InvoiceID"`
	Discounts       []Discount       `json:"discounts"  gorm:"foreignKey:InvoiceID"`
	Surcharges      []Surcharge      `json:"surcharges"  gorm:"foreignKey:InvoiceID"`
	SubTotal        float64          `json:"sub_total"`
	TotalDiscounts  float64          `json:"total_discounts,omitempty"`
	TotalSurcharges float64          `json:"total_surcharges,omitempty"`
	Tips            float64          `json:"tips"`
	BaseTax         float64          `json:"base_tax"`
	Taxes           float64          `json:"taxes"`
	Total           float64          `json:"total"`
	PaymentID       *uint            `json:"payment_id"`
	Payment         *payment.Payment `json:"payment"`
	CreatedAt       *time.Time       `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       *time.Time       `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt  `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Item struct {
	ID          uint            `json:"id"`
	InvoiceID   *uint           `json:"invoice_id"`
	ProductID   *uint           `json:"product_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	SKU         string          `json:"sku"`
	Price       float64         `json:"price" gorm:"precision:18;scale:2"`
	Comments    string          `json:"comments"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Discount struct {
	ID          uint           `json:"id"`
	InvoiceID   *uint          `json:"invoice_id"`
	Name        string         `json:"name,omitempty"`
	Type        string         `json:"type"`
	Percentage  float64        `json:"percentage,omitempty" gorm:"precision:18;scale:2"`
	Amount      float64        `json:"amount,omitempty" gorm:"precision:18;scale:2"`
	Description string         `json:"description,omitempty"`
	Terms       string         `json:"terms,omitempty"`
	ChannelID   *uint          `json:"channel_id,omitempty"`
	StoreID     *uint          `json:"store_id,omitempty"`
	BrandID     *uint          `json:"brand_id,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Surcharge struct {
	ID          uint            `json:"id"`
	InvoiceID   *uint           `json:"invoice_id"`
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
