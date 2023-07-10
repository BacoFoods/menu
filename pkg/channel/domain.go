package channel

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest      string = "error bad request"
	ErrorFindingChannel  string = "error finding channel"
	ErrorCreatingChannel string = "error creating channel"
	ErrorGettingChannel  string = "error getting channel"
	ErrorUpdatingChannel string = "error updating channel"
	ErrorDeletingChannel string = "error deleting channel"
)

type Channel struct {
	ID           uint            `json:"id"`
	Name         string          `json:"name"`
	ShortName    string          `json:"short_name"`
	Enabled      bool            `json:"enabled"`
	ShippingCost float64         `json:"shipping_cost,omitempty"`
	CreatedAt    *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt    *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt    *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Channel) (*Channel, error)
	Find(map[string]string) ([]Channel, error)
	Get(string) (*Channel, error)
	Update(*Channel) (*Channel, error)
	Delete(string) (*Channel, error)
}
