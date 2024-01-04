package facturacion

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	LogRepository string = "pkg/facturacion/repository"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Update(config *FacturacionConfig) (*FacturacionConfig, error) {
	if err := r.db.Save(config).Error; err != nil {
		return nil, err
	}

	return config, nil
}

func (r *Repository) Create(config *FacturacionConfig) error {
	return r.db.Create(config).Error
}

func (r *Repository) FindByStoreAndType(storeID uint, docType string) (*FacturacionConfig, error) {
	var config FacturacionConfig
	if err := r.db.Where("store_id = ? AND document_type = ?", storeID, docType).First(&config).Error; err != nil {
		shared.LogError("error getting facturacion config", LogRepository, "FindByStoreAndType", err, storeID, docType)
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf(ErrorFacturacionConfigNotFound)
		}
		return nil, fmt.Errorf(ErrorFacturacionConfigGetting)
	}

	return &config, nil
}

func (r *Repository) FindByStoreAndTypeAndIncrement(storeID uint, docType string) (*FacturacionConfig, error) {
	var config *FacturacionConfig
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("store_id = ? AND document_type = ?", storeID, docType).
			First(&config).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}

			return err
		}

		config.LastNumber = config.LastNumber + 1

		return tx.Save(config).Error
	})

	return config, err
}

func (r *Repository) FindByStore(storeID uint) ([]FacturacionConfig, error) {
	var config []FacturacionConfig
	if err := r.db.Where("store_id = ?", storeID).Find(&config).Error; err != nil {
		return nil, err
	}

	return config, nil
}
