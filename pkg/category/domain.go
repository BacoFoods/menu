package category

import (
	"github.com/BacoFoods/menu/pkg/product"
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID          uint              `json:"id"`
	Image       string            `json:"image"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Products    []product.Product `json:"products" gorm:"many2many:categories_products"`
	CreatedAt   *time.Time        `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt    `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type CategoriesProducts struct {
	ID         uint `json:"id"`
	CategoryID uint `json:"category_id" gorm:"primaryKey"`
	ProductID  uint `json:"product_id" gorm:"primaryKey"`
}
