package product

import (
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/taxes"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest      string = "error bad request"
	ErrorCreatingProduct string = "error creating product"
	ErrorFindingProduct  string = "error finding product"
	ErrorGettingProduct  string = "error getting product"
	ErrorUpdatingProduct string = "error updating product"
	ErrorDeletingProduct string = "error deleting product"
)

type Product struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Image       string             `json:"image"`
	SKU         string             `json:"sku"`
	Price       float32            `json:"price" gorm:"precision:18;scale:2"`
	TaxID       *uint              `json:"tax_id"`
	Tax         *taxes.Tax         `json:"tax" swaggerignore:"true"`
	DiscountID  *uint              `json:"discount_id"`
	Discount    *discount.Discount `json:"discount" gorm:"foreignKey:DiscountID" swaggerignore:"true"`
	Unit        string             `json:"unit"`
	BrandID     *uint              `json:"brand_id" binding:"required"`
	CreatedAt   *time.Time         `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt     `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Product) (*Product, error)
	Find(map[string]string) ([]Product, error)
	Get(productID []string) ([]Product, error)
	Update(*Product) (*Product, error)
	Delete(string) (*Product, error)
}
