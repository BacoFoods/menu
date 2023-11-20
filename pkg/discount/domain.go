package discount

import (
	"time"

	"gorm.io/gorm"
)

const (
	ErrorBadRequest       string = "error bad request"
	ErrorCreatingDiscount string = "error creating discount"
	ErrorGettingDiscount  string = "error getting discount"
	ErrorUpdatingDiscount string = "error updating discount"
	ErrorDeletingDiscount string = "error deleting discount"
	ErrorFindingDiscount  string = "error finding discount"

	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeValue      DiscountType = "value"
)

type DiscountType string

type Discount struct {
	ID          uint           `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Type        DiscountType   `json:"type"`
	Percentage  float64        `json:"percentage,omitempty" gorm:"precision:18;scale:2"`
	Value       float64        `json:"value,omitempty" gorm:"precision:18;scale:2"`
	Description string         `json:"description,omitempty"`
	Terms       string         `json:"terms,omitempty"`
	ChannelID   *uint          `json:"channel_id,omitempty"`
	StoreID     *uint          `json:"store_id,omitempty"`
	BrandID     *uint          `json:"brand_id,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type RepositoryI interface {
	Create(*Discount) (*Discount, error)
	Find(filters map[string]string) ([]Discount, error)
	Get(string) (*Discount, error)
	Update(Discount) (*Discount, error)
	Delete(string) (*Discount, error)
}
