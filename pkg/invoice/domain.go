package invoice

import (
	"math"
	"time"

	"github.com/BacoFoods/menu/pkg/payment"
	"github.com/BacoFoods/menu/pkg/tables"
	"gorm.io/gorm"
)

const (
	ErrorBadRequest                = "error bad request"
	ErrorInvoiceCreation           = "error creating invoice"
	ErrorGettingInvoice            = "error getting invoice"
	ErrorInvoiceUpdate             = "error updating invoice"
	ErrorInvalidTipAmount          = "error invalid tip amount"
	ErrorTipPercentageExceedsLimit = "error tip percentage exceeds limit"

	TaxPercentage = 0.08
)

type Repository interface {
	Create(invoice *Invoice) (*Invoice, error)
	Get(invoiceID string) (*Invoice, error)
	Update(invoice *Invoice) (*Invoice, error)
	UpdateTip(invoice *Invoice) (*Invoice, error)
}

type Invoice struct {
	ID              uint             `json:"id"`
	OrderID         *uint            `json:"order_id"`
	BrandID         *uint            `json:"brand_id" binding:"required"`
	StoreID         *uint            `json:"store_id" binding:"required"`
	ChannelID       *uint            `json:"channel_id" binding:"required"`
	TableID         *uint            `json:"table_id"`
	Table           *tables.Table    `json:"table"`
	Items           []Item           `json:"items"  gorm:"foreignKey:InvoiceID"`
	Discounts       []Discount       `json:"discounts" gorm:"many2many:invoice_discounts;"`
	Surcharges      []Surcharge      `json:"surcharges" gorm:"many2many:invoice_surcharges;"`
	SubTotal        float64          `json:"sub_total"`
	TotalDiscounts  float64          `json:"total_discounts,omitempty"`
	TotalSurcharges float64          `json:"total_surcharges,omitempty"`
	Tips            float64          `json:"tips"`
	Type            string           `json:"type"`
	BaseTax         float64          `json:"base_tax"`
	Taxes           float64          `json:"taxes"`
	Total           float64          `json:"total"`
	PaymentID       *uint            `json:"payment_id"`
	Payment         *payment.Payment `json:"payment"`
	CreatedAt       *time.Time       `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       *time.Time       `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt  `json:"deleted_at,omitempty" swaggerignore:"true"`
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
	Invoices    []Invoice      `json:"invoices" gorm:"many2many:invoice_discounts;"`
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
	Invoices    []Invoice       `json:"invoices" gorm:"many2many:invoice_surcharges;"`
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

type InvoiceDiscount struct {
	gorm.Model
	InvoiceID  uint
	DiscountID uint
}

type InvoiceSurcharge struct {
	gorm.Model
	InvoiceID   uint
	SurchargeID uint
}

// SetupDiscountInvoicesJoinTable configures a many-to-many relationship between Discount and Invoice
// This function should be called during migration to set up the join table invoice_discounts.
// For more information, refer to the GORM documentation: https://gorm.io/docs/many_to_many.html#Customize-JoinTable
func SetupDiscountInvoicesJoinTable(db *gorm.DB) error {
	return db.SetupJoinTable(&Discount{}, "Invoices", &InvoiceDiscount{})
}

// SetupSurchargeInvoicesJoinTable configures a many-to-many relationship between Surcharge and Invoice
// This function should be called during migration to set up the join table invoice_surcharges.
// For more information, refer to the GORM documentation: https://gorm.io/docs/many_to_many.html#Customize-JoinTable
func SetupSurchargeInvoicesJoinTable(db *gorm.DB) error {
	return db.SetupJoinTable(&Surcharge{}, "Invoices", &InvoiceSurcharge{})
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
