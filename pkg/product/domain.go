package product

import (
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/taxes"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest       string = "error bad request"
	ErrorCreatingProduct  string = "error creating product"
	ErrorFindingProduct   string = "error finding product"
	ErrorGettingProduct   string = "error getting product"
	ErrorUpdatingProduct  string = "error updating product"
	ErrorDeletingProduct  string = "error deleting product"
	ErrorAddingModifier   string = "error adding modifier"
	ErrorRemovingModifier string = "error removing modifier"

	ErrorModifierCreation        string = "error creating modifier"
	ErrorModifierAddingProduct   string = "error adding product to modifier"
	ErrorModifierRemovingProduct string = "error removing product from modifier"
	ErrorModifierGetting         string = "error getting modifiers"
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
	Modifiers   []Modifier         `json:"modifiers" gorm:"many2many:product_modifiers;"`
	CreatedAt   *time.Time         `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt    `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Modifier struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	Products    []Product       `json:"products" swaggerignore:"true" gorm:"many2many:modifier_products;"`
	BrandID     *uint           `json:"brand_id" binding:"required"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Product) (*Product, error)
	Find(map[string]string) ([]Product, error)
	Get(productID string) (*Product, error)
	GetByIDs(productIDs []string) ([]Product, error)
	Update(*Product) (*Product, error)
	Delete(string) (*Product, error)
	AddModifier(product *Product, modifier *Modifier) (*Product, error)
	RemoveModifier(product *Product, modifier *Modifier) (*Product, error)
	GetOverriders(productID, field string) ([]Overrider, error)

	ModifierCreate(*Modifier) (*Modifier, error)
	ModifierGet(modifierID string) (*Modifier, error)
	ModifierFind(map[string]string) ([]Modifier, error)
	ModifierAddProduct(product *Product, modifier *Modifier) (*Modifier, error)
	ModifierRemoveProduct(product *Product, modifier *Modifier) (*Modifier, error)
}

type Entity struct {
	Code  string
	Label string
}

var Entities map[string]Entity = map[string]Entity{
	"name": {
		Code:  "name",
		Label: "Nombre",
	},
	"description": {
		Code:  "description",
		Label: "Descripción",
	},
	"image": {
		Code:  "image",
		Label: "Imagen",
	},
	"price": {
		Code:  "price",
		Label: "Precio",
	},
	"enable": {
		Code:  "enable",
		Label: "Habilitado",
	},
}
