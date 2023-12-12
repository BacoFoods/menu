package invoice

import (
	"fmt"
	"time"

	"strings"

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
	if invoice.ID == 0 {
		if err := r.db.Create(invoice).Error; err != nil {
			shared.LogError("error creating invoice", LogRepository, "CreateUpdate", err, *invoice)
			return nil, err
		}
	}

	var invoiceDB Invoice
	if err := r.db.First(&invoiceDB, invoice.ID).Error; err != nil {
		shared.LogError("error getting invoice", LogRepository, "CreateUpdate", err, invoice.ID, *invoice)
		return nil, err
	}

	if err := r.db.Model(&invoiceDB).Where("id = ?", invoice.ID).Updates(invoice).Error; err != nil {
		shared.LogError("error updating invoice", LogRepository, "CreateUpdate", err, invoice.ID, invoice)
		return nil, err

	}
	return &invoiceDB, nil
}

// Get method for get an invoice in database
func (r *DBRepository) Get(invoiceID string) (*Invoice, error) {
	if strings.TrimSpace(invoiceID) == "" {
		err := fmt.Errorf(ErrorInvoiceIDEmpty)
		shared.LogWarn("error getting invoice", LogRepository, "Get", err)
		return nil, err
	}

	var invoice Invoice

	if err := r.db.Preload(clause.Associations).
		First(&invoice, invoiceID).Error; err != nil {
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

// FindInvoices method for finding the most recent invoice for each order in the database with additional filters including date range
func (r *DBRepository) FindInvoices(filter map[string]interface{}) ([]Invoice, error) {

	tx := r.db.
		Preload(clause.Associations).
		Select("DISTINCT ON (order_id) *").
		Order("order_id, created_at DESC")

	// Handle filter for StoreID
	if storeID, ok := filter["store_id"]; ok {
		tx = tx.Where("store_id = ?", storeID)
		delete(filter, "store_id")
	}

	// Handle filter for multiple stores
	if storeIDs, ok := filter["stores"]; ok {
		tx = tx.Where("store_id IN ?", storeIDs)
		delete(filter, "stores")
	}

	// Handle date range filter for start_date
	handleDateRangeFilter(tx, filter, "start_date", "created_at >= ?")

	// Handle date range filter for end_date
	handleDateRangeFilter(tx, filter, "end_date", "created_at < ?") // Use < to include the entire specified end date

	// Execute the query and handle errors
	var invoices []Invoice
	if err := tx.Find(&invoices, filter).Error; err != nil {
		shared.LogError("error finding invoices", LogRepository, "FindInvoices", err, filter)
		return nil, err
	}

	return invoices, nil
}

func handleDateRangeFilter(tx *gorm.DB, filter map[string]interface{}, key, condition string) {
	if value, ok := filter[key]; ok {
		switch key {
		case "start_date", "end_date":
			timestamp, err := time.Parse("2006-01-02", value.(string))
			if err != nil {
				fmt.Println("Error parsing timestamp:", err)
				return
			}

			if key == "end_date" {
				timestamp = timestamp.Add(24*time.Hour - time.Second)
			}

			tx = tx.Where(condition, timestamp)
		default:
			fmt.Println("Unsupported filter:", key)
		}

		delete(filter, key)
	}
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
	if strings.TrimSpace(invoiceID) == "" {
		err := fmt.Errorf(ErrorInvoiceIDEmpty)
		shared.LogWarn("error deleting invoice", LogRepository, "Delete", err)
		return err
	}
	if err := r.db.Delete(&Invoice{}, invoiceID).Error; err != nil {
		shared.LogError("error deleting invoice", LogRepository, "Delete", err, invoiceID)
		return err
	}
	return nil
}

// Print to get a printable invoice from database
func (r *DBRepository) Print(invoiceID string) (*DTOPrintable, error) {
	if strings.TrimSpace(invoiceID) == "" {
		err := fmt.Errorf(ErrorInvoiceIDEmpty)
		shared.LogWarn("error printing invoice", LogRepository, "Print", err)
		return nil, err
	}

	var invoice DTOPrintable

	if err := r.db.Table("invoices as i").
		Select("s.name as store_name, s.address as store_address, s.phone as store_phone, b.name as brand_name, b.document as brand_document, b.city as brand_city, i.created_at as date, i.waiter, i.shift_id, c.name as client_name, c.document as client_document, c.email as client_email, c.address as client_address, o.id as order_id, t.display_name as table_name, i.sub_total as subtotal, i.total_discounts as discount, i.tip, i.tip_amount, i.total_surcharges as surcharge, i.base_tax, i.taxes").
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

	var items []Item
	if err := r.db.Find(&items, "invoice_id = ?", invoiceID).Error; err != nil {
		shared.LogError("error printing invoice", LogRepository, "Print", err, invoiceID)
		return nil, err
	}

	return &invoice, nil
}

// FindDiscountApplied to find invoices with discount applied
func (r *DBRepository) FindDiscountApplied() ([]DiscountApplied, error) {
	var invoices []DiscountApplied
	if err := r.db.Find(&invoices).Error; err != nil {
		shared.LogError("error finding invoices", LogRepository, "FindDiscountApplied", err)
		return nil, err
	}
	return invoices, nil
}

// RemoveDiscountApplied to remove a discount applied
func (r *DBRepository) RemoveDiscountApplied(discountAppliedID string) (DiscountApplied, error) {
	var discountApplied DiscountApplied
	if err := r.db.First(&discountApplied, discountAppliedID).Error; err != nil {
		shared.LogError("error getting discount applied", LogRepository, "RemoveDiscountApplied", err, discountAppliedID)
		return discountApplied, err
	}
	if err := r.db.Delete(&discountApplied).Error; err != nil {
		shared.LogError("error removing discount applied", LogRepository, "RemoveDiscountApplied", err, discountAppliedID)
		return discountApplied, err
	}
	return discountApplied, nil
}
