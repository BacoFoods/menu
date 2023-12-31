package payment

import (
	"time"

	"gorm.io/gorm"
)

const (
	ErrorPaymentCreating = "error creating payment"
	ErrorPaymentGetting  = "error getting payment"
	ErrorPaymentFinding  = "error finding payment"
	ErrorPaymentUpdating = "error updating payment"
	ErrorPaymentDeleting = "error deleting payment"
	ErrorPaymentIDEmpty  = "error payment id empty"

	ErrorPaymentMethodFinding       = "error finding payment method"
	ErrorPaymentMethodWrongCode     = "error payment method wrong code"
	ErrorPaymentMethodEmptyCode     = "error payment method empty code"
	ErrorPaymentMethodAlreadyExists = "error payment method already exists"
	ErrorPaymentMethodCreation      = "error payment method creation"

	PaymentStatusEmmited  = "emitted" // for electronic invoices emission
	PaymentStatusPaid     = "paid"
	PaymentStatusPending  = "pending"
	PaymentStatusCanceled = "canceled"

	PaymentMethodCash        = "cash"
	PaymentMethodCardVisa    = "card-visa"
	PaymentMethodCardMaster  = "card-master"
	PaymentMethodCardAmex    = "card-amex"
	PaymentMethodCardDinners = "card-dinners"
	PaymentMethodBold        = "bold"
	PaymentMethodBono        = "bono"
	PaymentMethodYuno        = "yuno"
	PaymentMethodUntracked   = "untracked"
)

type Repository interface {
	Get(paymentID string) (*Payment, error)
	Find(filter map[string]any) ([]Payment, error)
	Create(payment *Payment) (*Payment, error)
	Update(payment *Payment) (*Payment, error)
	Delete(paymentID string) (*Payment, error)

	FindPaymentMethods(filter map[string]any) ([]PaymentMethod, error)
	GetPaymentMethod(paymentMethodID string) (*PaymentMethod, error)
	CreatePaymentMethod(*PaymentMethod) (*PaymentMethod, error)
	CreateDefaultPaymentMethods(brandID *uint) ([]PaymentMethod, error)
	UpdatePaymentMethod(*PaymentMethod) (*PaymentMethod, error)
	DeletePaymentMethod(string) (*PaymentMethod, error)
}

type Payment struct {
	ID uint `json:"id"`
	// An invoice can have multiple payments
	InvoiceID *uint `json:"invoice_id" binding:"required"`

	// Payment method used.
	// This can come from manual input in the POS such as <code>cash</code>, <code>card</code>, <code>check</code>, etc.
	// or from order in table with <code>paylot</code> or <code>yuno</code>
	Method          string         `json:"method" binding:"required"`
	PaymentMethodID *uint          `json:"payment_method_id" binding:"required"`
	PaymentMethod   *PaymentMethod `json:"payment_method,omitempty" gorm:"foreignKey:PaymentMethodID"`

	// Quantity is the amount of money paid
	Quantity float64 `json:"quantity" gorm:"precision:18;scale:4" binding:"required"`

	// Tip is the amount of money paid
	Tip float64 `json:"tip" gorm:"precision:18;scale:4" binding:"required"`

	// TotalValue is the paid = quantity + tip
	TotalValue float64 `json:"total_value" gorm:"precision:18;scale:4" binding:"required"`

	// Code is the reference number of the payment
	Code        string          `json:"code"`
	Status      string          `json:"status" binding:"required"`
	Reference   string          `json:"reference"`
	CheckoutURL *string         `json:"checkout_url"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type PaymentMethod struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	Name        string          `json:"name"`
	BrandID     *uint           `json:"brand_id"`
	StoreID     *uint           `json:"store_id"`
	ChannelID   *uint           `json:"channel_id"`
	ShortName   string          `json:"short_name"`
	Code        string          `json:"code" binding:"required"`
	Description string          `json:"description"`
	PlemsiCode  string          `json:"plemsi_code"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
