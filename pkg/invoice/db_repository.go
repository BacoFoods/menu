package invoice

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogRepository = "pkg/invoice/repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) DBRepository {
	return DBRepository{db}
}

func (r DBRepository) Create(invoice *Invoice) (*Invoice, error) {
	if err := r.db.Create(invoice).Error; err != nil {
		shared.LogError("Error creating invoice", LogRepository, "Create", err, *invoice)
		return nil, err
	}

	return invoice, nil
}
