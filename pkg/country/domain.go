package country

import (
	"time"

	"github.com/BacoFoods/menu/pkg/currency"
	"gorm.io/gorm"
)

const (
	ErrorCountryBadRequest string = "error bad request"
	ErrorCountryCreating   string = "error creating country"
	ErrorCountryGetting    string = "error getting country"
	ErrorCountryUpdating   string = "error updating country"
	ErrorCountryDeleting   string = "error deleting country"
	ErrorCountryIDEmpty    string = "country id is empty"
)

type CountryISO string

const (
	AR CountryISO = "AR"
	CO CountryISO = "CO"
	MX CountryISO = "MX"
	PE CountryISO = "PE"
)

type Country struct {
	ID         uint               `json:"id"`
	Name       string             `json:"name"`
	ISOCode    CountryISO         `json:"iso_code,omitempty"`
	CurrencyID *uint              `json:"currency_id"`
	Currency   *currency.Currency `json:"currency,omitempty" gorm:"foreignKey:CurrencyID"`
	PhoneCode  string             `json:"phone_code,omitempty"`
	CreatedAt  *time.Time         `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt  *time.Time         `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt  gorm.DeletedAt     `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Country) (*Country, error)
	Find(map[string]string) ([]Country, error)
	Get(string) (*Country, error)
	Update(Country) (*Country, error)
	Delete(string) (*Country, error)
}
