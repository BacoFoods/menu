package course

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogRepository = "pkg/course/repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r DBRepository) Find(filter map[string]any) ([]Course, error) {
	var courses []Course
	if err := r.db.Find(&courses, filter).Error; err != nil {
		shared.LogError("error finding courses", LogRepository, "Find", err, filter)
		return nil, err
	}

	return courses, nil
}

func (r DBRepository) Get(courseID string) (Course, error) {
	var course Course
	if err := r.db.First(&course, courseID).Error; err != nil {
		shared.LogError("error finding course", LogRepository, "Get", err, courseID)
		return Course{}, err
	}

	return course, nil
}

func (r DBRepository) Create(course Course) (Course, error) {
	if err := r.db.Save(&course).Error; err != nil {
		shared.LogError("error creating course", LogRepository, "Create", err, course)
		return Course{}, err
	}

	return course, nil
}

func (r DBRepository) Delete(courseID string) (Course, error) {
	course, err := r.Get(courseID)
	if err != nil {
		shared.LogError("error getting course", LogRepository, "Delete", err, courseID)
		return Course{}, err
	}

	if err := r.db.Delete(&course).Error; err != nil {
		shared.LogError("error deleting course", LogRepository, "Delete", err, course)
		return Course{}, err
	}

	return course, nil
}
