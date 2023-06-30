package product

import (
	"github.com/BacoFoods/menu/pkg/taxes"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	SKU         string         `json:"sku"`
	Price       float32        `json:"price" gorm:"precision:18;scale:2"`
	TaxID       *uint          `json:"tax_id"`
	Tax         *taxes.Tax     `json:"tax"`
	Discount    float32        `json:"discount" gorm:"precision:18;scale:2"`
	Unit        string         `json:"unit"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
