package store

import (
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/spot"
)

const (
	ErrorCreatingStore = "error creating store"
	ErrorGettingStore  = "error getting store"
	ErrorUpdatingStore = "error updating store"
	ErrorDeletingStore = "error deleting store"
	ErrorFindingStore  = "error finding store"
	ErrorBadRequest    = "error bad request"
)

type Store struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Menus       []menu.Menu `json:"menus,omitempty" gorm:"foreignKey:StoreID"`
	Spots       []spot.Spot `json:"spots,omitempty" gorm:"foreignKey:StoreID"`
	CreatedAt   string      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   string      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   string      `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Store) (*Store, error)
	Find(map[string]string) ([]Store, error)
	Get(string) (*Store, error)
	Update(*Store) (*Store, error)
	Delete(string) (*Store, error)
}
