package availability

import (
	"gorm.io/gorm"
	"time"
)

type Entity string

const (
	ErrorEnablingEntity = "error enabling entity"
	ErrorBadRequest     = "error bad request"

	EntityMenu Entity = "menu"
)

type Availability struct {
	ID        uint            `json:"id"`
	Entity    string          `json:"entity"`
	EntityID  *uint           `json:"entity_id"`
	Enable    bool            `json:"enable"`
	Place     string          `json:"place"`
	PlaceID   *uint           `json:"place_id"`
	CreatedAt *time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`
}

type Repository interface {
	EnableEntity(entity, place string, entityID, placeID uint, enable bool) error
	FindEntityByPlace(Entity, string, string) ([]Availability, error)
}
