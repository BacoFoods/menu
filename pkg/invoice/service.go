package invoice

import "fmt"

const LogService = "pkg/invoice/service"

type Service interface {
	Get(invoiceID string) (*Invoice, error)
	Update(invoice *Invoice) (*Invoice, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}

// Get returns a single Invoice object by ID.
func (s service) Get(invoiceID string) (*Invoice, error) {
	return s.repository.Get(invoiceID)
}

// Update actualiza el Invoice con los nuevos datos proporcionados en updateData.
func (s service) Update(updateData *Invoice) (*Invoice, error) {
	// get an invoice by ID
	existingInvoice, err := s.repository.Get(fmt.Sprintf("%d", updateData.ID))
	if err != nil {
		return nil, err
	}

	// Update the fields with updateData values
	existingInvoice.Type = updateData.Type
	existingInvoice.Tips = updateData.Tips

	// Update PaymentID if exists
	if updateData.PaymentID != nil {
		existingInvoice.PaymentID = updateData.PaymentID
	}

	// Update SurchargeID if exists
	if len(updateData.Surcharges) > 0 {
		existingInvoice.Surcharges[0].ID = updateData.Surcharges[0].ID
	}

	// Update SurchargeID if exists
	if len(updateData.Discounts) > 0 {
		existingInvoice.Discounts[0].ID = updateData.Discounts[0].ID
	}

	// Update the invoice in DB
	updatedInvoice, err := s.repository.Update(existingInvoice)
	if err != nil {
		return nil, err
	}

	return updatedInvoice, nil
}
