package tables

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest    = "error bad request"
	ErrorTableUpdating = "error updating table"
	ErrorTableDeleting = "error deleting table"
	ErrorTableCreating = "error creating table"
	ErrorTableGetting  = "error getting table"
	ErrorTableFinding  = "error finding table"
)

type Table struct {
	ID          *uint           `json:"id,omitempty"`
	DisplayID   string          `json:"display_id" binding:"required"`
	DisplayName string          `json:"display_name" binding:"required"`
	Number      int             `json:"number" binding:"required"`
	XLocation   float32         `json:"xlocation,omitempty"`
	YLocation   float32         `json:"ylocation,omitempty"`
	IsActive    bool            `json:"is_active"`
	ZoneID      *uint           `json:"zone_id"`
	OrderID     *uint           `json:"order_id,omitempty"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Get(id string) (*Table, error)
	Find(query map[string]any) ([]Table, error)
	Create(table *Table) (*Table, error)
	Update(id string, table *Table) (*Table, error)
	Delete(id string) error
}
