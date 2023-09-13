package invoice

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
)

const LogRepository = "pkg/invoice/repository"

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db}
}

func (r DBRepository) CreateUpdate(invoice *Invoice) (*Invoice, error) {
	if err := r.db.Save(invoice).Error; err != nil {
		shared.LogError("Error creating invoice", LogRepository, "Create", err, *invoice)
		return nil, err
	}

	return invoice, nil
}

// Get method for get an invoice in database
func (r *DBRepository) Get(invoiceID string) (*Invoice, error) {
	var invoice Invoice

	if err := r.db.First(&invoice, invoiceID).Error; err != nil {
		shared.LogError("error getting invoice", LogRepository, "Get", err, invoiceID)
		return nil, err
	}

	return &invoice, nil
}
