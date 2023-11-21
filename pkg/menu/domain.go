package menu

import (
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/product"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorMenuBadRequest           string = "error bad request"
	ErrorMenuIDEmpty              string = "error menu id empty"
	ErrorMenuFinding              string = "error finding menu"
	ErrorMenuFindingByPlace       string = "error finding menu by place"
	ErrorMenuGetting              string = "error getting menu"
	ErrorMenuCreating             string = "error creating menu"
	ErrorMenuUpdating             string = "error updating menu"
	ErrorMenuDeleting             string = "error deleting menu"
	ErrorMenuUpdatingAvailability string = "error updating availability"
	ErrorMenuFindingChannels      string = "error finding channels"
	ErrorMenuAddingCategory       string = "error adding category"
	ErrorMenuRemovingCategory     string = "error removing category"
	ErrorMenuWrongBrand           string = "error adding category to menu wrong brand"
)

type Repository interface {
	Create(*Menu) (*Menu, error)
	Find(map[string]string) ([]Menu, error)
	Get(string) (*Menu, error)
	Update(*Menu) (*Menu, error)
	Delete(string) (*Menu, error)
	FindByPlace(string, string) ([]Menu, error)
	GetMenuItems(string) ([]Item, error)
	AddCategory(menuID string, category *category.Category) (*Menu, error)
	RemoveCategory(menuID string, category *category.Category) (*Menu, error)
}

type Menu struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Categories  []category.Category `json:"categories" gorm:"many2many:menus_categories"`
	StartTime   *time.Time          `json:"start_time,omitempty"`
	EndTime     *time.Time          `json:"end_time,omitempty"`
	BrandID     *uint               `json:"brand_id"`
	Enable      bool                `json:"enable"`
	CreatedAt   *time.Time          `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt      `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type MenusCategories struct {
	MenuID     *uint      `json:"menu_id" gorm:"primaryKey"`
	CategoryID *uint      `json:"category_id" gorm:"primaryKey"`
	Enable     bool       `json:"enable,omitempty"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" swaggerignore:"true"`
	// DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

// JoinTable this function allows to generate a many2many relation between entities
// this function is called by migration to associate the tables to this relation
// please take a look to MenusCategories struct, this struct is used
// to generate the table menus_categories with more fields like enable, start_time, end_time, etc.
// More information in gorm docs: https://gorm.io/docs/many_to_many.html#Customize-JoinTable
func (b Menu) JoinTable(db gorm.DB) error {
	return db.SetupJoinTable(&Menu{}, "Categories", &MenusCategories{})
}

type Item struct {
	product.Product
	CategoryID    *uint  `json:"category_id"`
	OverriderName string `json:"overrider"`
}

var precedence = map[string]int{
	"brand":   3,
	"store":   2,
	"channel": 1,
}

func IsAllowOverride(item Item, overrider product.Overrider) bool {
	return precedence[item.OverriderName] < precedence[overrider.Name]
}

func OverrideProducts(items []Item, overriders []product.Overrider) map[uint][]product.Product {
	itemsByCategories := make(map[uint][]product.Product, 0)

	for _, item := range items {
		var prod product.Product

		// If product doesn't have overriders, then use the product as it is.
		if len(overriders) == 0 {
			prod = product.Product{
				ID:          item.ID,
				Name:        item.Name,
				Description: item.Description,
				Image:       item.Image,
				SKU:         item.SKU,
				Price:       item.Price,
				TaxID:       item.TaxID,
				DiscountID:  item.DiscountID,
				Unit:        item.Unit,
			}
			// If product has overriders, then use the overrider product.
		} else {
			for _, overrider := range overriders {
				if item.ID == *overrider.ProductID && IsAllowOverride(item, overrider) {
					prod = product.Product{
						ID:          item.ID,
						Name:        item.Name,
						Description: overrider.Description,
						Image:       overrider.Image,
						SKU:         item.SKU,
						Price:       overrider.Price,
						TaxID:       item.TaxID,
						DiscountID:  item.DiscountID,
						Unit:        item.Unit,
					}
				}
			}
		}

		itemsByCategories[*item.CategoryID] = append(itemsByCategories[*item.CategoryID], prod)
	}

	return itemsByCategories
}
