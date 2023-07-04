package taxes

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogDBRepository string = "pkg/tax/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(tax *Tax) (*Tax, error) {

	if err := r.db.Create(tax).Error; err != nil {
		shared.LogError("error creating tax", LogDBRepository, "Create", err, tax)
		return nil, err
	}
	return tax, nil
}

func (r *DBRepository) Find(query map[string]string) ([]Tax, error) {
	var tax []Tax
	if err := r.db.Preload(clause.Associations).Find(&tax, query).Error; err != nil {
		shared.LogError("error finding tax", LogDBRepository, "Find", err)
		return nil, err
	}
	return tax, nil
}

func (r *DBRepository) Get(taxID string) (*Tax, error) {
	var tax Tax
	if err := r.db.Preload(clause.Associations).First(&tax, taxID).Error; err != nil {
		shared.LogError("error getting tax", LogDBRepository, "Find", err, taxID)
		return nil, err
	}
	return &tax, nil
}

func (r *DBRepository) Update(tax Tax) (*Tax, error) {
	var taxDB Tax
	if err := r.db.First(&taxDB, tax.ID).Error; err != nil {
		shared.LogError("error getting tax", LogDBRepository, "Update", err, tax.ID, tax)
		return nil, err
	}

	if err := r.db.Model(&taxDB).Updates(tax).Error; err != nil {
		shared.LogError("error updating tax", LogDBRepository, "Update", err, tax.ID, tax, tax)
		return nil, err
	}
	return &taxDB, nil
}

func (r *DBRepository) Delete(taxID string) (*Tax, error) {
	var taxDB Tax
	if err := r.db.First(&taxDB, taxID).Error; err != nil {
		shared.LogError("error getting tax", LogDBRepository, "Delete", err, taxID)
		return nil, err
	}

	if err := r.db.Delete(&taxDB, taxID).Error; err != nil {
		shared.LogError("error deleting tax", LogDBRepository, "Delete", err, taxDB)
		return nil, err
	}
	return &taxDB, nil
}
