package shift

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const (
	LogDBRepository string = "pkg/shift/db_repository"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(shift *Shift) (*Shift, error) {
	if err := r.db.Create(shift).Error; err != nil {
		shared.LogError("failed to create new shift", LogDBRepository, "Create", err)
		return nil, err
	}

	return shift, nil
}

func (r *DBRepository) Update(shift *Shift) (*Shift, error) {
	if err := r.db.Save(shift).Error; err != nil {
		shared.LogError("failed to update shift", LogDBRepository, "Update", err)
		return nil, err
	}

	return shift, nil
}

func (r *DBRepository) GetOpenShift(storeID *uint) (*Shift, error) {
	var shift Shift
	if err := r.db.Where("store_id = ? AND end_time IS NULL", storeID).First(&shift).Error; err != nil {
		shared.LogError("failed to get open shift", LogDBRepository, "GetOpenShift", err)
		return nil, err
	}

	return &shift, nil
}
