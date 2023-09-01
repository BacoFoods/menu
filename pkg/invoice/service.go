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
	// Obtener el invoice existente por ID
	existingInvoice, err := s.repository.Get(fmt.Sprintf("%d", updateData.ID))
	if err != nil {
		return nil, err
	}

	// Actualiza los campos con los valores proporcionados en updateData
	existingInvoice.Type = updateData.Type
	existingInvoice.Tips = updateData.Tips
	// Actualiza otros campos según sea necesario

	// Actualiza PaymentID si existe
	if updateData.PaymentID != nil {
		existingInvoice.PaymentID = updateData.PaymentID
	}

	// Actualiza SurchargeID si existe
	if len(updateData.Surcharges) > 0 {
		existingInvoice.Surcharges[0].ID = updateData.Surcharges[0].ID
	}

	// Actualiza DiscountID si existe
	if len(updateData.Discounts) > 0 {
		existingInvoice.Discounts[0].ID = updateData.Discounts[0].ID
	}

	// Actualizar el invoice en la base de datos
	updatedInvoice, err := s.repository.Update(existingInvoice)
	if err != nil {
		return nil, err
	}

	return updatedInvoice, nil
}
