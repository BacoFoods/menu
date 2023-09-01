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
	// Obtener el invoice
	invoice, err := s.repository.Get(fmt.Sprintf("%d", updateData.ID))
	if err != nil {
		return nil, err
	}

	// Actualiza los campos con los valores proporcionados en updateData
	invoice.Type = updateData.Type
	invoice.Tips = updateData.Tips
	// Actualiza otros campos según sea necesario

	// Actualiza PaymentID si existe
	if updateData.PaymentID != nil {
		invoice.PaymentID = updateData.PaymentID
	}

	// Actualiza SurchargeID si existe
	if len(updateData.Surcharges) > 0 {
		invoice.Surcharges[0].ID = updateData.Surcharges[0].ID
	}

	// Actualiza DiscountID si existe
	if len(updateData.Discounts) > 0 {
		invoice.Discounts[0].ID = updateData.Discounts[0].ID
	}

	// Actualizar el invoice en la base de datos
	updatedInvoice, err := s.repository.Update(invoice)
	if err != nil {
		return nil, err
	}

	return updatedInvoice, nil
}
