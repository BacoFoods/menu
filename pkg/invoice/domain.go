package invoice

import (
	"fmt"
	"math"
	"time"

	"github.com/BacoFoods/menu/pkg/client"

	"github.com/BacoFoods/menu/pkg/payment"
	"gorm.io/gorm"
)

const (
	ErrorBadRequest                      = "error bad request"
	ErrorInvoiceCreation                 = "error creating invoice"
	ErrorGettingInvoice                  = "error getting invoice"
	ErrorInvoiceFinding                  = "error finding invoices"
	ErrorInvoiceUpdate                   = "error updating invoice"
	ErrorInvalidTipAmount                = "error invalid tip amount"
	ErrorTipPercentageExceedsLimit       = "error tip percentage exceeds limit"
	ErrorTipPercentageValue              = "error tip percentage wrong value"
	ErrorInvoiceAddingClient             = "error adding client to invoice"
	ErrorInvoiceRemovingClient           = "error removing client from invoice"
	ErrorInvoiceWrongClient              = "error wrong client for invoice"
	ErrorItemNotFound                    = "error item not found"
	ErrorInvoiceSeparatingNotEnoughItems = "error separating invoice not enough items sent"

	TaxPercentage     = 0.08
	TipTypePercentage = "PERCENTAGE"
	TipTypeAmount     = "AMOUNT"
	TipPercentageMax  = 0.1
)

type Repository interface {
	CreateUpdate(invoice *Invoice) (*Invoice, error)
	Get(invoiceID string) (*Invoice, error)
	Find(filter map[string]any) ([]Invoice, error)
	UpdateTip(invoice *Invoice) (*Invoice, error)
	CreateBatch(invoices []Invoice) ([]Invoice, error)
	Delete(invoiceID string) error
}

type Invoice struct {
	ID              uint              `json:"id"`
	OrderID         *uint             `json:"order_id"`
	BrandID         *uint             `json:"brand_id" binding:"required"`
	StoreID         *uint             `json:"store_id" binding:"required"`
	ChannelID       *uint             `json:"channel_id" binding:"required"`
	TableID         *uint             `json:"table_id"`
	Items           []Item            `json:"items"  gorm:"foreignKey:InvoiceID"`
	Discounts       []Discount        `json:"discounts"  gorm:"foreignKey:InvoiceID"`
	Surcharges      []Surcharge       `json:"surcharges"  gorm:"foreignKey:InvoiceID"`
	SubTotal        float64           `json:"sub_total"`
	TotalDiscounts  float64           `json:"total_discounts,omitempty"`
	TotalSurcharges float64           `json:"total_surcharges,omitempty"`
	Tips            float64           `json:"tips"`
	BaseTax         float64           `json:"base_tax"`
	Taxes           float64           `json:"taxes"`
	Total           float64           `json:"total"`
	Payments        []payment.Payment `json:"payments" gorm:"foreignKey:InvoiceID"`
	ClientID        *uint             `json:"client_id"`
	Client          *client.Client    `json:"client,omitempty"`
	CreatedAt       *time.Time        `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       *time.Time        `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt   `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (i *Invoice) CalculateTip(value float64, tipType string) error {
	i.BaseTax = math.Ceil(i.SubTotal / (1 + TaxPercentage))
	i.Taxes = i.SubTotal - i.BaseTax

	switch tipType {
	case TipTypePercentage:
		if value > TipPercentageMax {
			return fmt.Errorf(ErrorTipPercentageExceedsLimit)
		} else if !(value == 0.05) && !(value == 0.1) {
			return fmt.Errorf(ErrorTipPercentageValue)
		}
		i.Tips = math.Floor(value * i.BaseTax)
	case TipTypeAmount:
		if value < 0 {
			return fmt.Errorf(ErrorInvalidTipAmount)
		}
		i.Tips = math.Floor(value + i.BaseTax)
	}

	i.Total = i.BaseTax + i.Taxes + i.TotalSurcharges - i.TotalDiscounts + i.Tips
	return nil
}

func (i *Invoice) MapItems() map[uint]Item {
	items := make(map[uint]Item)
	for _, item := range i.Items {
		items[item.ID] = item
	}
	return items
}

type Item struct {
	ID          uint            `json:"id"`
	InvoiceID   *uint           `json:"invoice_id"`
	ProductID   *uint           `json:"product_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	SKU         string          `json:"sku"`
	Price       float64         `json:"price" gorm:"precision:18;scale:2"`
	Comments    string          `json:"comments"`
	Hash        string          `json:"hash"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Discount struct {
	ID          uint           `json:"id"`
	InvoiceID   *uint          `json:"invoice_id"`
	Name        string         `json:"name,omitempty"`
	Type        string         `json:"type"`
	Percentage  float64        `json:"percentage,omitempty" gorm:"precision:18;scale:2"`
	Amount      float64        `json:"amount,omitempty" gorm:"precision:18;scale:2"`
	Description string         `json:"description,omitempty"`
	Terms       string         `json:"terms,omitempty"`
	ChannelID   *uint          `json:"channel_id,omitempty"`
	StoreID     *uint          `json:"store_id,omitempty"`
	BrandID     *uint          `json:"brand_id,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Surcharge struct {
	ID          uint            `json:"id"`
	InvoiceID   *uint           `json:"invoice_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Percentage  float64         `json:"percentage" gorm:"precision:18;scale:2"`
	Amount      float64         `json:"amount" gorm:"precision:18;scale:2"`
	Active      bool            `json:"active"`
	ChannelID   *uint           `json:"channel_id,omitempty"`
	StoreID     *uint           `json:"store_id,omitempty"`
	BrandID     *uint           `json:"brand_id,omitempty"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
