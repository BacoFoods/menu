package store

import (
	"time"

	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/tables"
	"gorm.io/gorm"
)

const (
	ErrorStoreCreation              = "error creating store"
	ErrorStoreGet                   = "error getting store"
	ErrorStoreUpdate                = "error updating store"
	ErrorStoreDelete                = "error deleting store"
	ErrorStoreFind                  = "error finding store"
	ErrorBadRequest                 = "error bad request"
	ErrorStoreEnable                = "error enabling store"
	ErrorStoreAddingChannel         = "error adding channel"
	ErrorStoreZonesGettingByStoreID = "error getting zones by store id"
	ErrorStoreGettingChannels       = "error getting channels by brand id"
)

type Store struct {
	ID        uint              `json:"id"`
	Code      string            `json:"code"`
	Name      string            `json:"name"`
	BrandID   *uint             `json:"brand_id"`
	Enabled   bool              `json:"enabled"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Image     string            `json:"image"`
	City      string            `json:"city"`
	Channels  []channel.Channel `json:"channels,omitempty" gorm:"many2many:store_channels;" swaggerignore:"true"`
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
	Enable(storeID string) (*Store, error)
	AddChannel(storeID string, channel *channel.Channel) (*Store, error)

	FindZonesByStore(storeID string) ([]tables.Zone, error)
	GetZoneByStore(storeID, zoneID string) (*tables.Zone, error)
}
