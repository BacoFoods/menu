package invoice

import "github.com/BacoFoods/menu/pkg/payment"

const LogService = "pkg/invoice/service"

type Service interface {
	Get(invoiceID string) (*Invoice, error)
	Update(invoiceID string, updateData map[string]interface{}) (*Invoice, error)
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
func (s service) Update(invoiceID string, updateData map[string]interface{}) (*Invoice, error) {
	// Obtener el invoice
	invoice, err := s.repository.Get(invoiceID)
	if err != nil {
		return nil, err
	}

	// Actualiza los campos con los valores proporcionados en updateData
	for field, value := range updateData {
		switch field {
		case "Type":
			if typeValue, ok := value.(string); ok {
				invoice.Type = typeValue
			}
		case "Payment":
			if paymentValue, ok := value.(*payment.Payment); ok {
				invoice.Payment = paymentValue
			}
		case "Surcharges":
			if surchargesValue, ok := value.([]Surcharge); ok {
				invoice.Surcharges = surchargesValue
			}
		case "Tips":
			if tipsValue, ok := value.(float64); ok {
				invoice.Tips = tipsValue
			}
		case "Discounts":
			if discountsValue, ok := value.([]Discount); ok {
				invoice.Discounts = discountsValue
			}
		// Agrega más casos para otros campos que quieras actualizar
		}
	}

	// Actualizar el invoice en la base de datos
	updatedInvoice, err := s.repository.Update(invoice)
	if err != nil {
		return nil, err
	}

	return updatedInvoice, nil
}
