package category

import (
	"github.com/BacoFoods/menu/pkg/product"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorCategoryBadRequest      string = "error bad request"
	ErrorCategoryIDEmpty         string = "error category id empty"
	ErrorCategoryFinding         string = "error finding category"
	ErrorCategoryGetting         string = "error getting category"
	ErrorCategoryCreating        string = "error creating category"
	ErrorCategoryUpdating        string = "error updating category"
	ErrorCategoryDeleting        string = "error deleting category"
	ErrorCategoryGettingMenus    string = "error getting menus from category id"
	ErrorCategoryAddingProduct   string = "error adding product to category"
	ErrorCategoryRemovingProduct string = "error removing product from category"
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
