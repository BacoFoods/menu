package channel

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorChannelBadRequest string = "error bad request"
	ErrorChannelFinding    string = "error finding channel"
	ErrorChannelCreating   string = "error creating channel"
	ErrorChannelGetting    string = "error getting channel"
	ErrorChannelUpdating   string = "error updating channel"
	ErrorChannelDeleting   string = "error deleting channel"
	ErrorChannelIDEmpty    string = "error channel id empty"
)

type Channel struct {
	ID           uint            `json:"id"`
	Name         string          `json:"name"`
	ShortName    string          `json:"short_name"`
	Enabled      bool            `json:"enabled"`
	ShippingCost float64         `json:"shipping_cost,omitempty" gorm:"precision:18;scale:2"`
	BrandID      *uint           `json:"brand_id" binding:"required"`
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

	FindByIDs([]string) ([]Channel, error)
}
