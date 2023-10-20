package currency

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest       string = "error bad request"
	ErrorCreatingCurrency string = "error creating currency"
	ErrorGettingCurrency  string = "error getting currency"
	ErrorUpdatingCurrency string = "error updating currency"
	ErrorDeletingCurrency string = "error deleting currency"
)

type Currency struct {
	ID        uint           `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Code      string         `json:"code,omitempty"`
	Symbol    string         `json:"symbol,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Currency) (*Currency, error)
	Find(map[string]string) ([]Currency, error)
	Get(string) (*Currency, error)
	Update(Currency) (*Currency, error)
	Delete(string) (*Currency, error)
}
