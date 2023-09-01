package invoice

import (
	"fmt"
)

// Builder is a builder design pattern implementation for Invoice
type InvoiceBuilder struct {
	Errors     []error
	PaymentID  *uint
	Surcharge  Surcharge
	Discount   Discount
	Tips       float64
	Type       string
	SubTotal   float64
}

// NewInvoiceBuilder returns a new builder for creating an Invoice
func NewInvoiceBuilder() *InvoiceBuilder {
	return new(InvoiceBuilder)
}

// SetType sets the type for the Invoice
func (ib *InvoiceBuilder) SetType(typ string) *InvoiceBuilder {
	ib.Type = typ
	return ib
}

// SetPaymentID sets the payment ID for the Invoice
func (ib *InvoiceBuilder) SetPaymentID(paymentID uint) *InvoiceBuilder {
	ib.PaymentID = &paymentID
	return ib
}

// SetSurchargeID sets the surcharge ID for the Invoice
func (ib *InvoiceBuilder) SetSurchargeID(surchargeID uint) *InvoiceBuilder {
	ib.Surcharge.ID = surchargeID
	return ib
}

// SetTips sets the tips for the Invoice (with a business rule check)
func (ib *InvoiceBuilder) SetTips(tips float64) *InvoiceBuilder {
	// Check if tips exceed 10%
	if tips > 0.1*ib.SubTotal {
		ib.Errors = append(ib.Errors, fmt.Errorf("tips cannot exceed 10%% of the subtotal"))
	} else {
		ib.Tips = tips
	}
	return ib
}

// SetDiscountID sets the discount ID for the Invoice
func (ib *InvoiceBuilder) SetDiscountID(discountID uint) *InvoiceBuilder {
	ib.Discount.ID = discountID
	return ib
}

// Build returns a new Invoice and an error if there are any validation errors
func (ib *InvoiceBuilder) Build() (*Invoice, error) {
	// Validate your builder fields here
	// ...

	// Create a new instance of Invoice
	invoice := &Invoice{
		Type:        ib.Type,
		PaymentID:   ib.PaymentID,
		Surcharges:  []Surcharge{ib.Surcharge},
		Tips:        ib.Tips,
		Discounts:   []Discount{ib.Discount},
		SubTotal:    ib.SubTotal,
		// Other fields of Invoice...
	}

	return invoice, nil
}


