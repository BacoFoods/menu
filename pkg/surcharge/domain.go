package surcharge

import (
	"gorm.io/gorm"
	"time"
)

type Surcharge struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Percentage  float64         `json:"percentage" gorm:"precision:18;scale:2"`
	Amount      float64         `json:"amount" gorm:"precision:18;scale:2"`
	Active      bool            `json:"active"`
	ChannelID   *uint           `json:"channel_id,omitempty"`
	StoreID     *uint           `json:"store_id,omitempty"`
	BrandID     *uint           `json:"brand_id,omitempty"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Find(filters map[string]string) ([]Surcharge, error)
	Get(id string) (*Surcharge, error)
	Create(surcharge *Surcharge) (*Surcharge, error)
	Update(id string, surcharge *Surcharge) (*Surcharge, error)
	Delete(id string) (*Surcharge, error)
}
