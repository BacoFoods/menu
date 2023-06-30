package currency

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LogGormRepository string = "pkg/currency/repository"

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db}
}

func (r *GormRepository) Create(currency *Currency) (*Currency, error) {
	if err := r.db.Save(currency).Error; err != nil {
		shared.LogError("error creating currency", LogGormRepository, "Create", err, currency)
		return nil, err
	}
	return currency, nil
}

func (r *GormRepository) Find(query map[string]string) ([]Currency, error) {
	var currency []Currency
	if err := r.db.Preload(clause.Associations).Find(&currency, query).Error; err != nil {
		shared.LogError("error finding currency", LogGormRepository, "Find", err)
		return nil, err
	}
	return currency, nil
}

func (r *GormRepository) Get(currencyID string) (*Currency, error) {
	var currency Currency
	if err := r.db.Preload(clause.Associations).First(&currency, currencyID).Error; err != nil {
		shared.LogError("error getting currency", LogGormRepository, "Find", err, currencyID)
		return nil, err
	}
	return &currency, nil
}

func (r *GormRepository) Update(currency Currency) (*Currency, error) {
	var currencyDB Currency
	if err := r.db.First(&currencyDB, currency.ID).Error; err != nil {
		shared.LogError("error getting currency", LogGormRepository, "Update", err, currency.ID, currency)
		return nil, err
	}

	if err := r.db.Model(&currencyDB).Updates(currency).Error; err != nil {
		shared.LogError("error updating currency", LogGormRepository, "Update", err, currency.ID, currency, currency)
		return nil, err
	}
	return &currencyDB, nil
}

func (r *GormRepository) Delete(currencyID string) (*Currency, error) {
	var currencyDB Currency
	if err := r.db.First(&currencyDB, currencyID).Error; err != nil {
		shared.LogError("error getting currency", LogGormRepository, "Delete", err, currencyID)
		return nil, err
	}

	if err := r.db.Delete(&currencyDB, currencyID).Error; err != nil {
		shared.LogError("error deleting currency", LogGormRepository, "Delete", err, currencyDB)
		return nil, err
	}
	return &currencyDB, nil
}
