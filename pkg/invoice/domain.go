package invoice

import (
	"gorm.io/gorm"
	"time"
)

type Invoice struct {
	ID        uint            `json:"id"`
	OrderID   *uint           `json:"order_id"`
	CreatedAt *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
