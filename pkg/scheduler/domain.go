package scheduler

import (
	"github.com/BacoFoods/menu/pkg/store"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	ErrorScheduleFinding                  = "error finding schedules"
	ErrorScheduleCreating                 = "error creating schedule"
	ErrorScheduleUpdating                 = "error updating schedule"
	ErrorScheduleDeleting                 = "error deleting schedule"
	ErrorScheduleFindingTodayStore        = "error finding today schedule for store"
	ErrorScheduleFindingTodayBrand        = "error finding today schedule for brand"
	ErrorScheduleCreatingDayAlreadyExists = "error creating schedule, day already exists"
	ErrorHolidayCreating                  = "error creating holiday"
	ErrorHolidayUpdating                  = "error updating holiday"
	ErrorHolidayDeleting                  = "error deleting holiday"
	ErrorHolidayFinding                   = "error finding holiday"
)

type Schedule struct {
	ID        uint            `json:"id"`
	StoreID   *uint           `json:"store_id" binding:"required"`
	Store     *store.Store    `json:"store" gorm:"foreignKey:StoreID"`
	BrandID   *uint           `json:"brand_id" binding:"required"`
	Day       string          `json:"day" binding:"required" enums:"monday,tuesday,wednesday,thursday,friday,saturday,sunday,holiday"`
	Opening   string          `json:"open" binding:"required"`
	Closing   string          `json:"close" binding:"required"`
	Enable    bool            `json:"enable"`
	CreatedAt *time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (s *Schedule) IsOpen(holiday *Holiday) bool {
	if holiday != nil && holiday.Enable == true && s.Day == "holiday" && s.Enable {
		now := time.Now().UTC().Add(-5 * time.Hour)

		opening, _ := time.Parse("15:04", s.Opening)
		nowOpening := time.Date(now.Year(), now.Month(), now.Day(), opening.Hour(), opening.Minute(), opening.Second(), opening.Nanosecond(), time.UTC)

		closing, _ := time.Parse("15:04", s.Closing)
		nowClosing := time.Date(now.Year(), now.Month(), now.Day(), closing.Hour(), closing.Minute(), closing.Second(), closing.Nanosecond(), time.UTC)

		response := now.After(nowOpening) && now.Before(nowClosing)
		return response
	}

	// If the day is not the same as today, return false, so it doesn't show as open
	if strings.ToLower(time.Now().Weekday().String()) != s.Day {
		return false
	}

	if !s.Enable {
		return false
	}

	now := time.Now().UTC().Add(-5 * time.Hour)

	opening, _ := time.Parse("15:04", s.Opening)
	nowOpening := time.Date(now.Year(), now.Month(), now.Day(), opening.Hour(), opening.Minute(), opening.Second(), opening.Nanosecond(), time.UTC)

	closing, _ := time.Parse("15:04", s.Closing)
	nowClosing := time.Date(now.Year(), now.Month(), now.Day(), closing.Hour(), closing.Minute(), closing.Second(), closing.Nanosecond(), time.UTC)

	response := now.After(nowOpening) && now.Before(nowClosing)
	return response
}

func (s *Schedule) ToMap() map[string]any {
	return map[string]any{
		"store_id": *s.StoreID,
		"brand_id": *s.BrandID,
		"day":      s.Day,
		"opening":  s.Opening,
		"closing":  s.Closing,
		"enable":   s.Enable,
	}
}

type Repository interface {
	Find(filter map[string]any) ([]Schedule, error)
	Create(schedule *Schedule) error
	Update(schedule *Schedule) error
	Delete(schedule *Schedule) error
	TodayStore(storeID string) (*Schedule, error)
	TodayBrand(brandID string) ([]Schedule, error)
	EnableStore(storeID string, enable bool) ([]Schedule, error)
	CreateHoliday(*Holiday) (*Holiday, error)
	UpdateHoliday(*Holiday) (*Holiday, error)
	DeleteHoliday(*Holiday) error
	FindHoliday() ([]Holiday, error)
	GetTodayHoliday() (*Holiday, error)
}

type Holiday struct {
	ID        uint            `json:"id"`
	BrandID   *uint           `json:"brand_id" validate:"required"`
	Day       time.Time       `json:"day" validate:"required"`
	Enable    bool            `json:"enable"`
	CountryID *uint           `json:"country_id" validate:"required"`
	CreatedAt *time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
