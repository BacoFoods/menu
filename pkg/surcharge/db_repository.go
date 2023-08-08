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

func (r *DBRepository) Get(surchargeID string) (*Surcharge, error) {
	var surcharge Surcharge

	if err := r.db.First(&surcharge, surchargeID).Error; err != nil {
		shared.LogError("error getting surcharge", LogDBRepository, "Get", err, surchargeID)
		return nil, err
	}

	return &surcharge, nil
}

func (r *DBRepository) Create(surcharge *Surcharge) (*Surcharge, error) {
	if err := r.db.Save(surcharge).Error; err != nil {
		shared.LogError("error creating surcharge", LogDBRepository, "Create", err, surcharge)
		return nil, err
	}

	return surcharge, nil
}

func (r *DBRepository) Update(surchargeID string, surcharge *Surcharge) (*Surcharge, error) {
	var surchargeDB Surcharge

	if err := r.db.First(&surchargeDB, surchargeID).Error; err != nil {
		shared.LogError("error getting surcharge", LogDBRepository, "Update", err, surcharge)
		return nil, err
	}

	if err := r.db.Model(&surchargeDB).Updates(surcharge).Error; err != nil {
		shared.LogError("error updating surcharge", LogDBRepository, "Update", err, surcharge)
		return nil, err
	}

	return surcharge, nil
}

func (r *DBRepository) Delete(surchargeID string) (*Surcharge, error) {
	var surcharge Surcharge

	if err := r.db.First(&surcharge, surchargeID).Error; err != nil {
		shared.LogError("error getting surcharge", LogDBRepository, "Delete", err, surchargeID)
		return &surcharge, err
	}

	if err := r.db.Delete(&surcharge).Error; err != nil {
		shared.LogError("error deleting surcharge", LogDBRepository, "Delete", err, surchargeID)
		return &surcharge, err
	}

	return &surcharge, nil
}
