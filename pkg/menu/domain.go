package menu

import (
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/channel"
	"gorm.io/gorm"
	"time"
)

type Menu struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Categories  []category.Category `json:"categories" gorm:"many2many:menus_categories"`
	Channels    []channel.Channel   `json:"channels" gorm:"many2many:menus_channels"`
	StartTime   *time.Time          `json:"start_time,omitempty"`
	EndTime     *time.Time          `json:"end_time,omitempty"`
	CreatedAt   *time.Time          `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt      `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Menu) (*Menu, error)
	Find() ([]Menu, error)
	Get(string) (*Menu, error)
	Update(*Menu) (*Menu, error)
	Delete(string) (*Menu, error)
}

type MenusCategories struct {
	ID           uint           `json:"id"`
	MenuID       uint           `json:"menu_id" gorm:"primaryKey"`
	CategoryID   uint           `json:"category_id" gorm:"primaryKey"`
	Availability bool           `json:"availability,omitempty"`
	StartTime    *time.Time     `json:"start_time,omitempty"`
	EndTime      *time.Time     `json:"end_time,omitempty"`
	CreatedAt    *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type MenusChannels struct {
	ID           uint           `json:"id"`
	MenuID       uint           `json:"menu_id" gorm:"primaryKey"`
	ChannelID    uint           `json:"channel_id" gorm:"primaryKey"`
	Availability bool           `json:"availability,omitempty"`
	StartTime    *time.Time     `json:"start_time,omitempty"`
	EndTime      *time.Time     `json:"end_time,omitempty"`
	CreatedAt    *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
