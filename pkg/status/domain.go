package status

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest     string = "error bad request"
	ErrorCreatingStatus string = "error creating status"
	ErrorDeletingStatus string = "error deleting status"
	ErrorUpdatingStatus string = "error updating status"
	ErrorGettingStatus  string = "error getting status"
)

type Status struct {
	ID           uint           `json:"id,omitempty" gorm:"primaryKey"`
	Name         string         `json:"name"`
	Color        string         `json:"color"`
	Code         string         `json:"code"`
	Active       bool           `json:"active"`
	PrevStatusID *uint          `json:"prev_status_id"`
	NextStatusID *uint          `json:"next_status_id"`
	Prev         *Status        `json:"prev" gorm:"foreignKey:PrevStatusID"`
	Next         *Status        `json:"next" gorm:"foreignKey:NextStatusID"`
	CreatedAt    *time.Time     `json:"created_at,omitempty"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(status *Status) (*Status, error)
	Delete(status string) error
	Update(status *Status, statusID string) (*Status, error)
	Get(statusID string) (*Status, error)
	Find() ([]Status, error)
	GetFirst() (*Status, error)
	GetByCode(code string) (*Status, error)
}
