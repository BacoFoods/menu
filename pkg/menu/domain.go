package menu

import (
	"github.com/BacoFoods/menu/pkg/category"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorBadRequest   string = "error bad request"
	ErrorFindingMenu  string = "error finding menu"
	ErrorGettingMenu  string = "error getting menu"
	ErrorCreatingMenu string = "error creating menu"
	ErrorUpdatingMenu string = "error updating menu"
	ErrorDeletingMenu string = "error deleting menu"
)

type Menu struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Categories  []category.Category `json:"categories" gorm:"many2many:menus_categories"`
	StartTime   *time.Time          `json:"start_time,omitempty"`
	EndTime     *time.Time          `json:"end_time,omitempty"`
	CreatedAt   *time.Time          `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt      `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Repository interface {
	Create(*Menu) (*Menu, error)
	Find(map[string]string) ([]Menu, error)
	Get(string) (*Menu, error)
	Update(*Menu) (*Menu, error)
	Delete(string) (*Menu, error)
}

type MenusCategories struct {
	ID         uint           `json:"id"`
	MenuID     uint           `json:"menu_id" gorm:"primaryKey"`
	CategoryID uint           `json:"category_id" gorm:"primaryKey"`
	Enable     bool           `json:"enable,omitempty"`
	StartTime  *time.Time     `json:"start_time,omitempty"`
	EndTime    *time.Time     `json:"end_time,omitempty"`
	CreatedAt  *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt  *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
