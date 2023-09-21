package invoice

import (
	"math"
	"time"

	"github.com/BacoFoods/menu/pkg/client"

	"github.com/BacoFoods/menu/pkg/payment"
	"gorm.io/gorm"
)

const (
	ErrorBadRequest            		= "error bad request"
	ErrorInvoiceCreation       		= "error creating invoice"
	ErrorGettingInvoice        		= "error getting invoice"
	ErrorInvoiceFinding        		= "error finding invoices"
	ErrorInvoiceUpdate         		= "error updating invoice"
	ErrorInvalidTipAmount         	= "error invalid tip amount"
	ErrorTipPercentageExceedsLimit 	= "error tip percentage exceeds limit"
	ErrorInvoiceAddingClient   		= "error adding client to invoice"
	ErrorInvoiceRemovingClient 		= "error removing client from invoice"
	ErrorInvoiceWrongClient    		= "error wrong client for invoice"

	TaxPercentage = 0.08
)

type Repository interface {
	CreateUpdate(invoice *Invoice) (*Invoice, error)
	Get(invoiceID string) (*Invoice, error)
	Find(filter map[string]any) ([]Invoice, error)
	UpdateTip(invoice *Invoice) (*Invoice, error)
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

type Item struct {
	ID          uint            `json:"id"`
	InvoiceID   *uint           `json:"invoice_id"`
	ProductID   *uint           `json:"product_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	SKU         string          `json:"sku"`
	Price       float64         `json:"price" gorm:"precision:18;scale:2"`
	Comments    string          `json:"comments"`
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

// ReCalculateTips recalcula el campo 'tips' y actualiza el campo 'total' del Invoice.
func (i *Invoice) ReCalculateTips() {
	tipsAmount := 0.0

	i.BaseTax = math.Round(i.SubTotal / (1 + TaxPercentage))

	if i.Tips == 0.1 {
		tipsAmount = math.Round(i.BaseTax * 0.1)
		i.Tips = tipsAmount
	} else if i.Tips > 1.0 {
		tipsAmount = math.Round(i.Tips)
		i.Tips = tipsAmount
	} else {
		tipsAmount = 0.0
		i.Tips = tipsAmount
	}

	i.Total = math.Round(i.BaseTax + i.Taxes + i.TotalSurcharges - i.TotalDiscounts + tipsAmount)
}