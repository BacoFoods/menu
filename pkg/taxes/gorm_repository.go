package taxes

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogGormRepository string = "pkg/tax/gorm_repository"

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db}
}

func (r *GormRepository) Create(tax *Tax) (*Tax, error) {

	if err := r.db.Create(tax).Error; err != nil {
		shared.LogError("error creating tax", LogGormRepository, "Create", err, tax)
		return nil, err
	}
	return tax, nil
}

func (r *GormRepository) Find(query map[string]string) ([]Tax, error) {
	var tax []Tax
	if err := r.db.Preload(clause.Associations).Find(&tax, query).Error; err != nil {
		shared.LogError("error finding tax", LogGormRepository, "Find", err)
		return nil, err
	}
	return tax, nil
}

func (r *GormRepository) Get(taxID string) (*Tax, error) {
	var tax Tax
	if err := r.db.Preload(clause.Associations).First(&tax, taxID).Error; err != nil {
		shared.LogError("error getting tax", LogGormRepository, "Find", err, taxID)
		return nil, err
	}
	return &tax, nil
}

func (r *GormRepository) Update(tax Tax) (*Tax, error) {
	var taxDB Tax
	if err := r.db.First(&taxDB, tax.ID).Error; err != nil {
		shared.LogError("error getting tax", LogGormRepository, "Update", err, tax.ID, tax)
		return nil, err
	}

	if err := r.db.Model(&taxDB).Updates(tax).Error; err != nil {
		shared.LogError("error updating tax", LogGormRepository, "Update", err, tax.ID, tax, tax)
		return nil, err
	}
	return &taxDB, nil
}

func (r *GormRepository) Delete(taxID string) (*Tax, error) {
	var taxDB Tax
	if err := r.db.First(&taxDB, taxID).Error; err != nil {
		shared.LogError("error getting tax", LogGormRepository, "Delete", err, taxID)
		return nil, err
	}

	if err := r.db.Delete(&taxDB, taxID).Error; err != nil {
		shared.LogError("error deleting tax", LogGormRepository, "Delete", err, taxDB)
		return nil, err
	}
	return &taxDB, nil
}
