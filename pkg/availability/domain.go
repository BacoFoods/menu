package availability

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Entity string
type Place string

const (
	ErrorEnablingEntity      = "error enabling entity"
	ErrorBadRequest          = "error bad request"
	ErrorEntityNotFound      = "error entity not found"
	ErrorPlaceNotFound       = "error place not found"
	ErrorFindingAvailability = "error finding availability"
	ErrorGettingAvailability = "error getting availability"

	EntityMenu     Entity = "menu"
	EntityCategory Entity = "category"
	PlaceStore     Place  = "store"
	PlaceChannel   Place  = "channel"
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
	EnableEntity(entity Entity, place Place, entityID, placeID uint, enable bool) error
	FindEntityByPlace(entity Entity, place Place, placeID string) ([]Availability, error)
	FindPlacesByEntity(entity Entity, entityID uint, place string) ([]Availability, error)
	Get(entity Entity, place Place, entityID, placeID uint) (Availability, error)
	Find(entity Entity, place Place, entityID uint) ([]Availability, error)
}

func GetEntity(entity string) (Entity, error) {
	switch entity {
	case "menu":
		return EntityMenu, nil
	default:
		return "", fmt.Errorf(ErrorEntityNotFound)
	}
}

func GetPlace(place string) (Place, error) {
	switch place {
	case "store":
		return PlaceStore, nil
	case "channel":
		return PlaceChannel, nil
	default:
		return "", fmt.Errorf(ErrorPlaceNotFound)
	}
}

func MapAvailabilityByPlace(availabilities []Availability) map[uint]Availability {
	m := make(map[uint]Availability)
	for _, availability := range availabilities {
		m[*availability.PlaceID] = availability
	}
	return m
}
