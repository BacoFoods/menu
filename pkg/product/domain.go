package product

import (
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/BacoFoods/menu/pkg/taxes"
	"gorm.io/gorm"
	"strconv"
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
	ErrorGettingCategory  string = "error getting category"

	ErrorModifierCreation        string = "error creating modifier"
	ErrorModifierAddingProduct   string = "error adding product to modifier"
	ErrorModifierRemovingProduct string = "error removing product from modifier"
	ErrorModifierGetting         string = "error getting modifiers"
	ErrorModifierUpdate          string = "error updating modifier"

	ErrorOverriderCreating string = "error creating overriders"
	ErrorOverriderFinding  string = "error finding overriders"
	ErrorOverriderGetting  string = "error getting overriders"
	ErrorOverriderUpdating string = "error updating overriders"
	ErrorOverriderDeleting string = "error deleting overriders"

	LogDomain string = "pkg/product/domain"
)

type Product struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Image       string             `json:"image"`
	SKU         string             `json:"sku"`
	Price       float64            `json:"price" gorm:"precision:18;scale:2"`
	TaxID       *uint              `json:"tax_id"`
	Tax         *taxes.Tax         `json:"tax" swaggerignore:"true"`
	DiscountID  *uint              `json:"discount_id"`
	Discount    *discount.Discount `json:"discount" gorm:"foreignKey:DiscountID" swaggerignore:"true"`
	Unit        string             `json:"unit"`
	BrandID     *uint              `json:"brand_id" binding:"required"`
	Color       string             `json:"color"`
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
	ApplyPrice  float64         `json:"apply_price" gorm:"precision:18;scale:2"`
	Category    Category        `json:"category"`
	Products    []Product       `json:"products" swaggerignore:"true" gorm:"many2many:modifier_products;"`
	BrandID     *uint           `json:"brand_id" binding:"required"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Overrider struct {
	ID          uint               `json:"id"`
	ProductID   *uint              `json:"product_id"`
	Product     *Product           `json:"product,omitempty" gorm:"foreignKey:ProductID" swaggerignore:"true"`
	Place       string             `json:"place"`
	PlaceID     *uint              `json:"place_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Image       string             `json:"image"`
	Price       float64            `json:"price" gorm:"precision:18;scale:2"`
	Enable      bool               `json:"enable"`
	DiscountID  *uint              `json:"discount_id"`
	Discount    *discount.Discount `json:"discount" gorm:"foreignKey:DiscountID" swaggerignore:"true"`
	CreatedAt   *time.Time         `json:"created_at" swaggerignore:"true"`
	UpdatedAt   *time.Time         `json:"updated_at" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt    `json:"deleted_at" swaggerignore:"true"`
}

type Category string

type Repository interface {
	Create(*Product) (*Product, error)
	Find(map[string]string) ([]Product, error)
	Get(productID string) (*Product, error)
	GetByIDs(productIDs []string) ([]Product, error)
	Update(*Product) (*Product, error)
	Delete(string) (*Product, error)
	AddModifier(product *Product, modifier *Modifier) (*Product, error)
	RemoveModifier(product *Product, modifier *Modifier) (*Product, error)
	GetOverriders(productID, field string) ([]OverriderDTO, error)
	GetOverriderIDs(productID string) ([]uint, error)
	UpdateOverriders(ids []uint, field string, value any) error
	GetCategory(productID string) ([]CategoryDTO, error)

	ModifierCreate(*Modifier) (*Modifier, error)
	ModifierGet(modifierID string) (*Modifier, error)
	ModifierFind(map[string]string) ([]Modifier, error)
	ModifierAddProduct(product *Product, modifier *Modifier) (*Modifier, error)
	ModifierRemoveProduct(product *Product, modifier *Modifier) (*Modifier, error)
	ModifierUpdate(*Modifier) (*Modifier, error)

	OverriderCreate(*Overrider) (*Overrider, error)
	OverriderCreateAll([]Overrider) error
	OverriderFind(map[string]string) ([]Overrider, error)
	OverriderGet(string) (*Overrider, error)
	OverriderUpdate(*Overrider) (*Overrider, error)
	OverriderDelete(string) (*Overrider, error)
	OverriderFindByPlace(string, string) ([]Overrider, error)
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

func TransformValue(entity string, value string) any {
	switch entity {
	case "price":
		price, err := strconv.ParseFloat(value, 32)
		if err != nil {
			shared.LogError("error parsing price", LogDomain, "TransformValue", err)
			return nil
		}
		return price
	case "enable":
		enable, err := strconv.ParseBool(value)
		if err != nil {
			shared.LogError("error parsing enable", LogDomain, "TransformValue", err)
			return nil
		}
		return enable
	}
	return nil
}
