package spot

import (
	"github.com/BacoFoods/menu/pkg/channel"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorCreatingSpot = "error creating spot"
	ErrorGettingSpot  = "error getting spot"
	ErrorUpdatingSpot = "error updating spot"
	ErrorDeletingSpot = "error deleting spot"
	ErrorFindingSpot  = "error finding spot"
	ErrorBadRequest   = "error bad request"
)

type Spot struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	StoreID   uint              `json:"store_id"`
	Enabled   bool              `json:"enabled"`
	Channels  []channel.Channel `json:"channels,omitempty" gorm:"foreignKey:SpotID"`
	CreatedAt *time.Time        `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt *time.Time        `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt   `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Spot) (*Spot, error)
	Find(map[string]string) ([]Spot, error)
	Get(string) (*Spot, error)
	Update(*Spot) (*Spot, error)
	Delete(string) (*Spot, error)
}
