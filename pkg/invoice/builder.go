package invoice

import (
	"errors"
	"fmt"
)

// Builder is a builder design pattern implementation for Invoice
type InvoiceBuilder struct {
	Errors     []error
	PaymentID  *uint
	SurchargeID  *Surcharge
	DiscountID   *Discount
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
	ib.SurchargeID.ID = surchargeID
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
	ib.DiscountID.ID = discountID
	return ib
}

// Build returns a new Invoice and an error if there are any validation errors
func (ib *InvoiceBuilder) Build() (*Invoice, error) {
	var validationErrors []error

	if ib.Type == "" {
		validationErrors = append(validationErrors, errors.New("type is required"))
	}

	if ib.PaymentID == nil {
		validationErrors = append(validationErrors, errors.New("PaymentID is required"))
	}

	if ib.SurchargeID == nil {
		validationErrors = append(validationErrors, errors.New("surchargeID is required"))
	}

	if ib.Tips > 0.1*ib.SubTotal {
		validationErrors = append(validationErrors, errors.New("tips cannot exceed 10% of the subtotal"))
	}

	if ib.DiscountID == nil {
		validationErrors = append(validationErrors, errors.New("discountID is required"))
	}

	// Check if there are validation errors
	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation errors: %v", validationErrors)
	}

	// Create a new instance of Invoice
	invoice := &Invoice{
		Type:      ib.Type,
		PaymentID: ib.PaymentID,
		Tips:      ib.Tips,
		SubTotal:  ib.SubTotal,
	}

	if ib.SurchargeID != nil {
		invoice.Surcharges = []Surcharge{*ib.SurchargeID}
	}

	if ib.DiscountID != nil {
		invoice.Discounts = []Discount{*ib.DiscountID}
	}

	return invoice, nil
}