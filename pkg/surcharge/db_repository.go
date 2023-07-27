package surcharge

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogDBRepository string = "pkg/surcharge/db_repository"

type DBRepository struct {
	db *gorm.DB
}

// Endpoint listar Categories by ProductID

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Find(filters map[string]string) ([]Surcharge, error) {
	var surcharges []Surcharge

	if err := r.db.Find(&surcharges, filters).Error; err != nil {
		shared.LogError("error finding surcharges", LogDBRepository, "Find", err, filters)
		return nil, err
	}

	return surcharges, nil
}
