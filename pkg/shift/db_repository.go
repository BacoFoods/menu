package shift

import "gorm.io/gorm"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) OpenShift() error {
	return nil
}

func (r *DBRepository) CloseShift() error {
	return nil
}
