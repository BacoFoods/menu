package cashier

import "gorm.io/gorm"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Open() error {
	return nil
}

func (r *DBRepository) Close() error {
	return nil
}
