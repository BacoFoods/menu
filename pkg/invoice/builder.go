package invoice

import (
	"errors"
	"fmt"
)

// Builder is a builder design pattern implementation for Invoice
type Builder struct {
	Invoice
	Errors []error
}

// NewInvoiceBuilder returns a new builder for creating an Invoice
func NewInvoiceBuilder() *Builder {
	return &Builder{}
}

// SetType sets the type for the Invoice
func (ib *Builder) SetType(typ string) *Builder {
	if typ == "" {
		ib.Errors = append(ib.Errors, errors.New("type is required"))
	}
	ib.Type = typ
	return ib
}

// SetPaymentID sets the payment ID for the Invoice
func (ib *Builder) SetPaymentID(paymentID uint) *Builder {
	ib.PaymentID = &paymentID
	return ib
}

// SetTips sets the tips for the Invoice (with a business rule check)
func (ib *Builder) SetTips(tips float64) *Builder {
	if tips != 0 {
		ib.Tips = tips
	} else if tips < 0 {
		ib.Errors = append(ib.Errors, errors.New("tips must be greater than 0"))
	} else {
		ib.Tips = 0
	}
	return ib
}

// AddExistingDiscounts añade descuentos existentes a la Invoice en construcción
func (ib *Builder) AddExistingDiscounts(discounts []Discount) *Builder {
	ib.Discounts = append(ib.Discounts, discounts...)
	return ib
}

// AddExistingSurcharges añade recargos existentes a la Invoice en construcción
func (ib *Builder) AddExistingSurcharges(surcharges []Surcharge) *Builder {
	ib.Surcharges = append(ib.Surcharges, surcharges...)
	return ib
}

// Build returns a new Invoice and an error if there are any validation errors
func (ib *Builder) Build() (*Invoice, error) {
	if len(ib.Errors) > 0 {
		return nil, fmt.Errorf("validation errors: %v", ib.Errors)
	}
	return &ib.Invoice, nil
}
