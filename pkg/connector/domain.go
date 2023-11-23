package connector

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorEquivalenceBadRequest string = "error bad request"
	ErrorEquivalenceCreating   string = "error creating equivalence"
	ErrorEquivalenceGetting    string = "error getting equivalence"
	ErrorEquivalenceUpdating   string = "error updating equivalence"
	ErrorEquivalenceDeleting   string = "error deleting equivalence"
	ErrorEquivalenceIDEmpty    string = "equivalence id is empty"
	ErrorBadRequest            string = "error bad request"
	ErrorInternalServer        string = "internal server error"
)

type Equivalence struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	ChannelID string          `json:"channel_id"`
	ProductID string          `json:"product_id"`
	SiesaID   string          `json:"siesa_id"`
	CreatedAt *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Equivalence) (*Equivalence, error)
	FindReference(filter map[string]string) (*Equivalence, error)
}
