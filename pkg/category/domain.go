package category

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID          uint           `json:"id"`
	Image       string         `json:"image"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
