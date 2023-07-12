package overriders

import (
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/product"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorCreatingOverriders = "error creating overriders"
	ErrorFindingOverriders  = "error finding overriders"
	ErrorGettingOverriders  = "error getting overriders"
	ErrorUpdatingOverriders = "error updating overriders"
	ErrorDeletingOverriders = "error deleting overriders"
	ErrorBadRequest         = "error bad request"
)

type Overriders struct {
	ID          uint               `json:"id"`
	ProductID   *uint              `json:"product_id"`
	Product     *product.Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Place       string             `json:"place"`
	PlaceID     *uint              `json:"place_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Image       string             `json:"image"`
	Price       float32            `json:"price" gorm:"precision:18;scale:2"`
	Enable      bool               `json:"enable"`
	DiscountID  *uint              `json:"discount_id"`
	Discount    *discount.Discount `json:"discount" gorm:"foreignKey:DiscountID"`
	CreatedAt   *time.Time         `json:"created_at"`
	UpdatedAt   *time.Time         `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt    `json:"deleted_at"`
}

type Repository interface {
	Create(*Overriders) (*Overriders, error)
	Find(map[string]string) ([]Overriders, error)
	Get(string) (*Overriders, error)
	Update(*Overriders) (*Overriders, error)
	Delete(string) (*Overriders, error)
	FindByPlace(string, string) ([]Overriders, error)
}
