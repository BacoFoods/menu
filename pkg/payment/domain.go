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
)

type Repository interface {
	Get(paymentID string) (*Payment, error)
	Find(filter map[string]any) ([]Payment, error)
	Create(payment *Payment) (*Payment, error)
	Update(payment *Payment) (*Payment, error)
	Delete(paymentID string) (*Payment, error)
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
