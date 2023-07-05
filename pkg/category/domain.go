package category

import (
	"github.com/BacoFoods/menu/pkg/product"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest       string = "error bad request"
	ErrorFindingCategory  string = "error finding category"
	ErrorGettingCategory  string = "error getting category"
	ErrorCreatingCategory string = "error creating category"
	ErrorUpdatingCategory string = "error updating category"
	ErrorDeletingCategory string = "error deleting category"
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
	ID         uint  `json:"id"`
	CategoryID *uint `json:"category_id" gorm:"primaryKey"`
	ProductID  *uint `json:"product_id" gorm:"primaryKey"`
}

type Repository interface {
	Find(map[string]string) ([]Category, error)
	Get(string) (*Category, error)
	Create(*Category) (*Category, error)
	Update(*Category) (*Category, error)
	Delete(string) (*Category, error)
}
