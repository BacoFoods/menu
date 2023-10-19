package shift

import (
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	OpenShift() error
	CloseShift() error
}

type CashierShift struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	ChannelID uint            `json:"channel_id"`
	StoreID   uint            `json:"store_id"`
	BrandID   uint            `json:"brand_id"`
	StartTime *time.Time      `json:"start_time"`
	EndTime   *time.Time      `json:"end_time"`
	CreatedAt *time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
