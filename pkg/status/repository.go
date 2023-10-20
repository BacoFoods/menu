package status

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/status/repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(status *Status) (*Status, error) {
	if err := r.db.Save(status).Error; err != nil {
		shared.LogError("error creating status", LogDBRepository, "Create", err, *status)
		return nil, err
	}

	return status, nil
}

func (r *DBRepository) Delete(statusID string) error {
	if err := r.db.Delete(Status{}, statusID).Error; err != nil {
		shared.LogError("error deleting status", LogDBRepository, "Delete", err, statusID)
		return err
	}

	return nil
}

func (r *DBRepository) Update(status *Status, statusID string) (*Status, error) {
	var status_ Status
	if err := r.db.First(&status_, statusID).Error; err != nil {
		shared.LogError("error getting status", LogDBRepository, "Update", err, statusID, *status)
		return nil, err
	}

	if err := r.db.Model(&status_).Updates(status).Error; err != nil {
		shared.LogError("error updating status", LogDBRepository, "Update", err, statusID, status, *status)
		return nil, err
	}

	return &status_, nil
}

func (r *DBRepository) Get(statusID string) (*Status, error) {
	if statusID == "" {
		shared.LogWarn("error getting brand", LogDBRepository, "Get", shared.ErrorIDEmpty)
		return nil, shared.ErrorIDEmpty
	}

	var status Status
	if err := r.db.Preload(clause.Associations).First(&status, statusID).Error; err != nil {
		shared.LogError("error getting status", LogDBRepository, "Find", err, statusID)
		return nil, err
	}

	return &status, nil
}

func (r *DBRepository) Find() ([]Status, error) {
	var status []Status
	if err := r.db.Preload(clause.Associations).Find(&status).Error; err != nil {
		shared.LogError("error finding status", LogDBRepository, "Find", err)
		return nil, err
	}

	return status, nil
}

func (r *DBRepository) GetByCode(code string) (*Status, error) {
	var status Status
	if err := r.db.First(&status, "code = ?", code).Error; err != nil {
		shared.LogError("error getting status by code", LogDBRepository, "GetByCode", err, code)
		return nil, err
	}

	return &status, nil
}
