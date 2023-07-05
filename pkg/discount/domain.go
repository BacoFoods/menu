package discount

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest       string = "error bad request"
	ErrorCreatingDiscount string = "error creating discount"
	ErrorGettingDiscount  string = "error getting discount"
	ErrorUpdatingDiscount string = "error updating discount"
	ErrorDeletingDiscount string = "error deleting discount"
	ErrorFindingDiscount  string = "error finding discount"
)

type Discount struct {
	ID          uint           `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Percentage  float32        `json:"percentage,omitempty"`
	Description string         `json:"description,omitempty"`
	Terms       string         `json:"terms,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type RepositoryI interface {
	Create(*Discount) (*Discount, error)
	Find(map[string]string) ([]Discount, error)
	Get(string) (*Discount, error)
	Update(Discount) (*Discount, error)
	Delete(string) (*Discount, error)
}
