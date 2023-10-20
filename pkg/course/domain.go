package course

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorCourseFinding  string = "error finding courses"
	ErrorCourseGetting  string = "error getting course by id"
	ErrorCourseBinding  string = "error binding course"
	ErrorCourseCreating string = "error creating course"
	ErrorCourseDeleting string = "error deleting course"
)

type Repository interface {
	Find(map[string]any) ([]Course, error)
	Get(string) (Course, error)
	Create(Course) (Course, error)
	Delete(string) (Course, error)
}

type Course struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	ChannelID   *uint           `json:"channel_id"`
	StoreID     *uint           `json:"store_id"`
	BrandID     *uint           `json:"brand_id"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
