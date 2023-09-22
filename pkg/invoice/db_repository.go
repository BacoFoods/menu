package invoice

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	if invoiceID == "" {
		shared.LogWarn("error getting invoice", LogRepository, "Get", shared.ErrorIDEmpty)
		return nil, shared.ErrorIDEmpty
	}

	var invoice Invoice

	if err := r.db.Preload(clause.Associations).First(&invoice, invoiceID).Error; err != nil {
		shared.LogError("error getting invoice", LogRepository, "Get", err, invoiceID)
		return nil, err
	}

	return &invoice, nil
}

// Find method for find invoices in database
func (r *DBRepository) Find(filter map[string]interface{}) ([]Invoice, error) {
	var invoices []Invoice
	if err := r.db.Preload(clause.Associations).Find(&invoices, filter).Error; err != nil {
		shared.LogError("error finding invoices", LogRepository, "Find", err, filter)
		return nil, err
	}

	return invoices, nil
}

// UpdateTip update the field 'tips' of an Invoice in database.
func (r *DBRepository) UpdateTip(invoice *Invoice) (*Invoice, error) {
	var invoiceDB Invoice
	if err := r.db.First(&invoiceDB, invoice.ID).Error; err != nil {
		shared.LogError("error getting invoice", LogRepository, "UpdateTip", err, invoice.ID, invoice)
		return nil, err
	}
	if err := r.db.Model(&invoiceDB).Where("id = ?", invoice.ID).Updates(invoice).Error; err != nil {
		shared.LogError("error updating tip of an invoice", LogRepository, "UpdateTip", err, invoice.ID, invoice, invoice)
		return nil, err
	}
	return &invoiceDB, nil
}
