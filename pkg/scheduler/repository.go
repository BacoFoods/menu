package scheduler

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

const (
	LogRepository = "pkg/scheduler/repository"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) DBRepository {
	return DBRepository{db}
}

func (r DBRepository) Find(filter map[string]any) ([]Schedule, error) {
	var schedule []Schedule
	if err := r.db.Preload(clause.Associations).Find(&schedule, filter).Error; err != nil {
		shared.LogError("error finding schedule", LogRepository, "Get", err, filter)
		return nil, fmt.Errorf(ErrorScheduleFinding)
	}
	return schedule, nil
}

func (r DBRepository) Create(schedule *Schedule) error {
	var storeDaySchedule Schedule
	if err := r.db.Where("store_id = ? AND day = ?", schedule.StoreID, schedule.Day).First(&storeDaySchedule).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			shared.LogError("error finding schedule", LogRepository, "Create", err, *schedule)
			return fmt.Errorf(ErrorScheduleCreating)
		}
	}

	if storeDaySchedule.ID != 0 {
		shared.LogError("schedule already exists", LogRepository, "Create", nil, *schedule)
		return fmt.Errorf(ErrorScheduleCreatingDayAlreadyExists)
	}

	if err := r.db.Create(schedule).Error; err != nil {
		shared.LogError("error creating schedule", LogRepository, "Create", err, *schedule)
		return fmt.Errorf(ErrorScheduleCreating)
	}
	return nil
}

func (r DBRepository) Update(schedule *Schedule) error {
	if err := r.db.Save(schedule).Error; err != nil {
		shared.LogError("error updating schedule", LogRepository, "Update", err, *schedule)
		return fmt.Errorf(ErrorScheduleUpdating)
	}
	return nil
}

func (r DBRepository) Delete(schedule *Schedule) error {
	if err := r.db.Delete(schedule).Error; err != nil {
		shared.LogError("error deleting schedule", LogRepository, "Delete", err, *schedule)
		return fmt.Errorf(ErrorScheduleDeleting)
	}
	return nil
}

func (r DBRepository) TodayStore(storeID string) (*Schedule, error) {
	var schedule Schedule
	day := strings.ToLower(time.Now().Weekday().String())
	if err := r.db.Where("store_id = ? AND day = ? and enable = true", storeID, day).First(&schedule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		shared.LogError("error finding schedule", LogRepository, "Today", err, storeID)
		return &schedule, fmt.Errorf(ErrorScheduleFindingTodayStore)
	}
	return &schedule, nil
}

func (r DBRepository) TodayBrand(brandID string) ([]Schedule, error) {
	var schedules []Schedule
	day := strings.ToLower(time.Now().Weekday().String())
	if err := r.db.Where("brand_id = ? AND day = ?", brandID, day).Find(&schedules).Error; err != nil {
		shared.LogError("error finding schedules", LogRepository, "Today", err, brandID)
		return nil, fmt.Errorf(ErrorScheduleFindingTodayBrand)
	}
	return schedules, nil
}

func (r DBRepository) EnableStore(storeID string, enable bool) ([]Schedule, error) {
	if err := r.db.Model(Schedule{}).
		Where("store_id = ?", storeID).
		Updates(map[string]any{"enable": enable}).
		Error; err != nil {
		shared.LogError("error updating schedule", LogRepository, "EnableStore", err, storeID)
		return nil, fmt.Errorf(ErrorScheduleUpdating)
	}

	var schedules []Schedule
	if err := r.db.Where("store_id = ?", storeID).Find(&schedules).Error; err != nil {
		shared.LogError("error finding schedules", LogRepository, "EnableStore", err, storeID)
		return nil, fmt.Errorf(ErrorScheduleFinding)
	}

	return schedules, nil
}
