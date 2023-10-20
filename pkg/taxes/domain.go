package taxes

import (
	"github.com/BacoFoods/menu/pkg/country"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest  string = "error bad request"
	ErrorCreatingTax string = "error creating tax"
	ErrorGettingTax  string = "error getting tax"
	ErrorUpdatingTax string = "error updating tax"
	ErrorDeletingTax string = "error deleting tax"
)

type Tax struct {
	ID          uint             `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Percentage  float64          `json:"percentage,omitempty" gorm:"precision:18;scale:4;not null"`
	Description string           `json:"description,omitempty"`
	CountryID   *uint            `json:"country_id,omitempty"`
	Country     *country.Country `json:"country,omitempty" gorm:"foreignKey:CountryID"`
	CreatedAt   *time.Time       `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time       `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt   `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Tax) (*Tax, error)
	Find(map[string]string) ([]Tax, error)
	Get(string) (*Tax, error)
	Update(Tax) (*Tax, error)
	Delete(string) (*Tax, error)
}
