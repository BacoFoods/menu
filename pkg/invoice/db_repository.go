package invoice

import (
	"fmt"

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

func (r DBRepository) Create(invoice *Invoice) (*Invoice, error) {
	if err := r.db.Create(invoice).Error; err != nil {
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

// Update method for update an invoice in database
func (r *DBRepository) Update(invoice *Invoice) (*Invoice, error) {
	var invoiceDB Invoice

	fmt.Println("invoice",invoice)
	fmt.Println("invoice id",invoice.ID)
	fmt.Println("invoice DB",&invoiceDB)
	
	if err := r.db.Where("id = ?", invoice.ID).First(&invoiceDB, invoice.ID).Error; err != nil {
		shared.LogError("error getting invoice", LogRepository, "Update", err, invoice.ID, invoice)
		return nil, err
	}
	fmt.Println(r.db.Model(&invoiceDB))
	if err := r.db.Model(&invoiceDB).Select("tips", "type", "payment_id","discounts", "surcharges").Where("id = ?", invoice.ID).Updates(invoice).Error; err != nil {
		shared.LogError("error updating invoice", LogRepository, "Update", err, invoice.ID, invoice)
		return nil, err
	}
	return &invoiceDB, nil
}