package brand

import (
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/store"
)

const (
	ErrorBrandCreating   = "error creating brand"
	ErrorBrandGetting    = "error getting brand"
	ErrorBrandUpdating   = "error updating brand"
	ErrorBrandDeleting   = "error deleting brand"
	ErrorBrandFinding    = "error finding brand"
	ErrorBrandBadRequest = "error bad request"
	ErrorBrandIDEmpty    = "brand id is empty"
)

type Brand struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Document     string        `json:"document"`
	DocumentType string        `json:"document_type"`
	SocialName   string        `json:"social_name"`
	Menus        []menu.Menu   `json:"menus,omitempty" gorm:"foreignKey:BrandID"`
	Stores       []store.Store `json:"stores,omitempty" gorm:"foreignKey:BrandID"`
	City         string        `json:"city"`
	CreatedAt    string        `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt    string        `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt    string        `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Brand) (*Brand, error)
	Find(map[string]string) ([]Brand, error)
	Get(string) (*Brand, error)
	Update(*Brand) (*Brand, error)
	Delete(string) (*Brand, error)
}
