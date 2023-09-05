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

// AddSurcharge adds a surcharge to the Invoice
func (ib *Builder) AddSurcharge(surcharge Surcharge) *Builder {
	ib.Surcharges = append(ib.Surcharges, surcharge)
	return ib
}

// SetTips sets the tips for the Invoice (with a business rule check)
func (ib *Builder) SetTips(tips float64) *Builder {
	// Check if tips exceed 10% of the subtotal
	if tips > 0.1*ib.SubTotal {
		ib.Errors = append(ib.Errors, fmt.Errorf("tips cannot exceed 10%% of the subtotal"))
	} else {
		ib.Tips = tips
	}
	return ib
}

// AddDiscount adds a discount to the Invoice
func (ib *Builder) AddDiscount(discount Discount) *Builder {
	ib.Discounts = append(ib.Discounts, discount)
	return ib
}

// Build returns a new Invoice and an error if there are any validation errors
func (ib *Builder) Build() (*Invoice, error) {
	if len(ib.Errors) > 0 {
		return nil, fmt.Errorf("validation errors: %v", ib.Errors)
	}
	return &ib.Invoice, nil
}
