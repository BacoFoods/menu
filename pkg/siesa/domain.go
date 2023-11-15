package siesa

import (
	"time"

	"gorm.io/gorm"
)

const (
	ErrorBadRequest     = "error bad request"
	ErrorInternalServer = "internal server error"
)

type Reference struct {
	ID                       uint            `json:"id" gorm:"primaryKey"`
	Category                 string          `json:"categoria"`
	Popapp                   string          `json:"popapp"`
	ReferenciaPdv            string          `json:"referencia_pdv"`
	ReferenciaDeliveryInline string          `json:"referencia_delivery_inline"`
	RappiPickUp              string          `json:"rappi_pick_up"`
	RappiBacu                string          `json:"rappi_bacu"`
	DidiStu                  string          `json:"didi_stu"`
	DidiBacu                 string          `json:"didi_bacu"`
	BacuMarketplace          string          `json:"bacu_marketplace"`
	CreatedAt                *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt                *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt                *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Reference) error
	TruncateRecords() error
	Find(map[string]string) (*Reference, error)
}
