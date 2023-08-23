package brand

import (
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/store"
)

const (
	ErrorCreatingBrand = "error creating brand"
	ErrorGettingBrand  = "error getting brand"
	ErrorUpdatingBrand = "error updating brand"
	ErrorDeletingBrand = "error deleting brand"
	ErrorFindingBrand  = "error finding brand"
	ErrorBadRequest    = "error bad request"
)

type Brand struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	NIT         string        `json:"nit"`
	SocialName  string        `json:"social_name"`
	Menus       []menu.Menu   `json:"menus,omitempty" gorm:"foreignKey:BrandID"`
	Stores      []store.Store `json:"stores,omitempty" gorm:"foreignKey:BrandID"`
	CreatedAt   string        `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   string        `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   string        `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Brand) (*Brand, error)
	Find(map[string]string) ([]Brand, error)
	Get(string) (*Brand, error)
	Update(*Brand) (*Brand, error)
	Delete(string) (*Brand, error)
}
