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
	ErrorGettingMenus     string = "error getting menus from category id"
	ErrorAddingProduct    string = "error adding product to category"
	ErrorRemovingProduct  string = "error removing product from category"
)

type Category struct {
	ID          uint              `json:"id"`
	Image       string            `json:"image"`
	Name        string            `json:"name"`
	BrandID     *uint             `json:"brand_id" binding:"required"`
	Description string            `json:"description"`
	Enable      bool              `json:"enable"`
	Color       string            `json:"color"`
	SortID      int               `json:"sort_id"`
	Products    []product.Product `json:"products" gorm:"many2many:categories_products" swaggerignore:"true"`
	CreatedAt   *time.Time        `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt    `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type MenusCategory struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

type Repository interface {
	Find(map[string]string) ([]Category, error)
	Get(string) (*Category, error)
	Create(*Category) (*Category, error)
	Update(*Category) (*Category, error)
	Delete(string) (*Category, error)
	GetMenusByCategory(categoryID string) ([]MenusCategory, error)
	AddProduct(categoryID, productID uint) (*Category, error)
	RemoveProduct(categoryID, productID uint) (*Category, error)
}
