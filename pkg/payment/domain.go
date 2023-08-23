package payment

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	Method    string          `json:"method"`
	CreatedAt *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
