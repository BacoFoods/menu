package payment

import (
	"gorm.io/gorm"
	"time"
)

const (
	ErrorPaymentCreating = "error creating payment"
	ErrorPaymentGetting  = "error getting payment"
	ErrorPaymentFinding  = "error finding payment"
	ErrorPaymentUpdating = "error updating payment"
	ErrorPaymentDeleting = "error deleting payment"

	ErrorPaymentMethodFinding = "error finding payment method"
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
	UpdatePaymentMethod(*PaymentMethod) (*PaymentMethod, error)
	DeletePaymentMethod(string) (*PaymentMethod, error)
}

type Payment struct {
	ID        uint            `json:"id"`
	InvoiceID *uint           `json:"invoice_id" binding:"required"`
	Method    string          `json:"method" binding:"required"`
	Quantity  float32         `json:"quantity" gorm:"precision:18;scale:4" binding:"required"`
	Code      string          `json:"code"`
	CreatedAt *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type PaymentMethod struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	BrandID   *uint           `json:"brand_id"`
	StoreID   *uint           `json:"store_id"`
	ChannelID *uint           `json:"channel_id"`
	ShortName string          `json:"short_name"`
	Code      string          `json:"code"`
	Franchise string          `json:"franchise"`
	CreatedAt *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}
