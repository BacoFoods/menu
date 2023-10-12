package invoice

import (
	"fmt"
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
	tx := r.db.
		Preload(clause.Associations)

	if _, ok := filter["paid"]; ok {
		// TODO: not rules for paid
		delete(filter, "paid")
	}

	if closed, ok := filter["closed"]; ok {
		// closed define if the invoice has a payment
		if closed == "true" {
			tx = tx.Joins("JOIN payments ON payments.invoice_id = invoices.id").
				Where("payments.deleted_at IS NULL")
		}
		delete(filter, "closed")
	}

	if days, ok := filter["days"]; ok {
		tx.Where(fmt.Sprintf("created_at >= NOW() - INTERVAL '%s' DAY", days))
		shared.LogWarn("filtering by days", LogRepository, "Find", nil, filter)
		delete(filter, "days")
	}

	var invoices []Invoice
	if err := tx.Find(&invoices, filter).Error; err != nil {
		shared.LogError("error finding invoices", LogRepository, "Find", err, filter)
		return nil, err
	}

	return invoices, nil
}

// UpdateTip update the field 'tips' of an Invoice in database.
func (r *DBRepository) UpdateTip(invoice *Invoice) (*Invoice, error) {
	var invoiceDB Invoice
	if err := r.db.First(&invoiceDB, invoice.ID).Error; err != nil {
		shared.LogError("error getting invoice", LogRepository, "UpdateTip", err, invoice.ID, *invoice)
		return nil, err
	}
	if err := r.db.Model(&invoiceDB).Where("id = ?", invoice.ID).Updates(invoice).Error; err != nil {
		shared.LogError("error updating tip of an invoice", LogRepository, "UpdateTip", err, invoice.ID, invoice)
		return nil, err
	}
	return &invoiceDB, nil
}

// CreateBatch creates a batch of invoices in database.
func (r *DBRepository) CreateBatch(invoices []Invoice) ([]Invoice, error) {
	if err := r.db.Create(&invoices).Error; err != nil {
		shared.LogError("error creating batch of invoices", LogRepository, "CreateBatch", err, invoices)
		return nil, err
	}
	return invoices, nil
}

// Delete deletes an invoice in database.
func (r *DBRepository) Delete(invoiceID string) error {
	if invoiceID == "" {
		shared.LogWarn("error deleting invoice", LogRepository, "Delete", shared.ErrorIDEmpty)
		return shared.ErrorIDEmpty
	}
	if err := r.db.Delete(&Invoice{}, invoiceID).Error; err != nil {
		shared.LogError("error deleting invoice", LogRepository, "Delete", err, invoiceID)
		return err
	}
	return nil
}

// Print to get a printable invoice from database
func (r *DBRepository) Print(invoiceID string) (*DBDTOPrintInvoice, error) {
	if invoiceID == "" {
		shared.LogWarn("error printing invoice", LogRepository, "Print", shared.ErrorIDEmpty)
		return nil, shared.ErrorIDEmpty
	}

	var invoice DBDTOPrintInvoice

	if err := r.db.Table("invoices as i").
		Select("s.name as store_name, s.address as store_address, s.phone as store_phone, b.name as brand_name, b.document as brand_document, b.city as brand_city, i.created_at as date, i.waiter, i.cashier, c.name as client_name, c.document as client_document, c.email as client_email, c.address as client_address, o.id as order_id, t.display_name as table_name, i.sub_total as subtotal, i.total_discounts as discount, i.tip, i.tip_amount, i.total_surcharges as surcharge, i.base_tax, i.taxes").
		Joins("left join stores as s on i.store_id = s.id").
		Joins("left join brands as b on i.brand_id = b.id").
		Joins("left join orders as o on i.order_id = o.id").
		Joins("left join clients as c on i.client_id = c.id").
		Joins("left join tables as t on i.table_id = t.id").
		Where("i.deleted_at is null").
		Scan(&invoice).Error; err != nil {
		shared.LogError("error printing invoice", LogRepository, "Print", err, invoiceID)
		return nil, err
	}

	var items []DBDTOPrintInvoiceItem
	if err := r.db.Model(Item{}).Find(&items, "invoice_id = ?", invoiceID).Error; err != nil {
		shared.LogError("error printing invoice", LogRepository, "Print", err, invoiceID)
		return nil, err
	}

	return &invoice, nil
}
