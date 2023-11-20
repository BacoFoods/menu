package siesa

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/siesa/repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(reference *Reference) error {
	if err := r.db.Save(reference).Error; err != nil {
		shared.LogError("error creating reference", LogDBRepository, "Create", err, reference)
		return err
	}
	return nil
}

// TruncateRecords truncates all records in the YourModel table
func (r *DBRepository) TruncateRecords() error {
	if err := r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Reference{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *DBRepository) Find(filters map[string]string) (*Reference, error) {
	var reference Reference
	if err := r.db.Unscoped().Preload(clause.Associations).Find(&reference, filters).Limit(1).Error; err != nil {
		shared.LogError("error getting reference row", LogDBRepository, "Find", err, filters)
		return nil, err
	}
	return &reference, nil
}
