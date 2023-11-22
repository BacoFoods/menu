package currency

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

const LogDBRepository string = "pkg/currency/db_repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r *DBRepository) Create(currency *Currency) (*Currency, error) {
	if err := r.db.Save(currency).Error; err != nil {
		shared.LogError("error creating currency", LogDBRepository, "Create", err, currency)
		return nil, err
	}
	return currency, nil
}

func (r *DBRepository) Find(query map[string]string) ([]Currency, error) {
	var currency []Currency
	if err := r.db.Preload(clause.Associations).Find(&currency, query).Error; err != nil {
		shared.LogError("error finding currency", LogDBRepository, "Find", err)
		return nil, err
	}
	return currency, nil
}

func (r *DBRepository) Get(currencyID string) (*Currency, error) {
	if strings.TrimSpace(currencyID) == "" {
		err := fmt.Errorf(ErrorCurrencyIDEmpty)
		shared.LogWarn("error getting currency", LogDBRepository, "Get", err)
		return nil, err
	}

	var currency Currency
	if err := r.db.Preload(clause.Associations).First(&currency, currencyID).Error; err != nil {
		shared.LogError("error getting currency", LogDBRepository, "Find", err, currencyID)
		return nil, err
	}
	return &currency, nil
}

func (r *DBRepository) Update(currency Currency) (*Currency, error) {
	var currencyDB Currency
	if err := r.db.First(&currencyDB, currency.ID).Error; err != nil {
		shared.LogError("error getting currency", LogDBRepository, "Update", err, currency.ID, currency)
		return nil, err
	}

	if err := r.db.Model(&currencyDB).Updates(currency).Error; err != nil {
		shared.LogError("error updating currency", LogDBRepository, "Update", err, currency.ID, currency, currency)
		return nil, err
	}
	return &currencyDB, nil
}

func (r *DBRepository) Delete(currencyID string) (*Currency, error) {
	var currencyDB Currency
	if err := r.db.First(&currencyDB, currencyID).Error; err != nil {
		shared.LogError("error getting currency", LogDBRepository, "Delete", err, currencyID)
		return nil, err
	}

	if err := r.db.Delete(&currencyDB, currencyID).Error; err != nil {
		shared.LogError("error deleting currency", LogDBRepository, "Delete", err, currencyDB)
		return nil, err
	}
	return &currencyDB, nil
}
