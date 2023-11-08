package shift

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorGettingAccount   = "error getting account"
	ErrorGettingOpenShift = "error getting open shift"
)

type Repository interface {
	Create(*Shift) (*Shift, error)
	Update(*Shift) (*Shift, error)
	GetOpenShift(storeID *uint) (*Shift, error)
	GetLastShift(storeID string) (*Shift, error)
}

type Shift struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	StoreID      *uint           `json:"store_id"`
	BrandID      *uint           `json:"brand_id"`
	AccountID    *uint           `json:"account_id"`
	StartTime    *time.Time      `json:"start_time"`
	EndTime      *time.Time      `json:"end_time"`
	StartBalance float64         `json:"start_balance"`
	EndBalance   float64         `json:"end_balance"`
	CreatedAt    *time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt    *time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt    *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
