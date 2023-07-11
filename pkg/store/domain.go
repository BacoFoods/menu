package store

import (
	"github.com/BacoFoods/menu/pkg/channel"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorCreatingStore = "error creating store"
	ErrorGettingStore  = "error getting store"
	ErrorUpdatingStore = "error updating store"
	ErrorDeletingStore = "error deleting store"
	ErrorFindingStore  = "error finding store"
	ErrorBadRequest    = "error bad request"
	ErrorAddingChannel = "error adding channel"
)

type Store struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	BrandID   *uint             `json:"brand_id"`
	Enabled   bool              `json:"enabled"`
	Image     string            `json:"image"`
	Channels  []channel.Channel `json:"channels,omitempty" gorm:"many2many:store_channels;"`
	Latitude  float64           `json:"latitude"`
	Longitude float64           `json:"longitude"`
	Address   string            `json:"address"`
	CreatedAt *time.Time        `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt *time.Time        `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt   `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Store) (*Store, error)
	Find(map[string]string) ([]Store, error)
	Get(string) (*Store, error)
	Update(*Store) (*Store, error)
	Delete(string) (*Store, error)
	FindByIDs(storeIDs []string) ([]Store, error)
	AddChannel(storeID string, channel *channel.Channel) (*Store, error)
}
