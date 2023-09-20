package invoice

import "fmt"

const LogService = "pkg/invoice/service"

type Service interface {
	Get(invoiceID string) (*Invoice, error)
	Update(invoice *Invoice, discounts []uint, Surcharges []uint) (*Invoice, error)
	UpdateTip(invoice *Invoice) (*Invoice, error)
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
func (s service) Update(updateData *Invoice, discounts []uint, surcharges []uint) (*Invoice, error) {
	// get an invoice by ID
	existingInvoice, err := s.repository.Get(fmt.Sprintf("%d", updateData.ID))

	if updateData.Tips > 0.1*existingInvoice.SubTotal {
		return nil, fmt.Errorf("tips cannot be greater than 10 percent of subtotal")
	}

	// Agrega nuevos descuentos)
	if err != nil {
		return nil, err
	}
	// Actualiza los valores de Discounts según los IDs proporcionados
	updateData.Discounts = []Discount{}
	for _, discountID := range discounts {
		discount := Discount{
			ID: discountID,
		}
		updateData.Discounts = append(updateData.Discounts, discount)
	}

	// Actualiza los valores de Surcharges según los IDs proporcionados
	updateData.Surcharges = []Surcharge{}
	for _, surchargeID := range surcharges {
		surcharge := Surcharge{
			ID: surchargeID,
		}

		updateData.Surcharges = append(updateData.Surcharges, surcharge)
	}
	// Si no se encuentra la factura, devuelve una instancia vacía o una factura con valores predeterminados
	if existingInvoice == nil {
		return &Invoice{}, nil
	}
	// Update the invoice in DB
	updatedInvoice, err := s.repository.Update(updateData)
	if err != nil {
		return nil, err
	}

	return updatedInvoice, nil
}

// UpdateTip actualiza el campo 'tips' de un Invoice y verifica si es un valor válido.
func (s service) UpdateTip(updateData *Invoice) (*Invoice, error) {
	// Obtener el invoice existente por su ID
	existingInvoice, err := s.repository.Get(fmt.Sprintf("%d", updateData.ID))
	if err != nil {
		return nil, err
	}

	// Verificar si las propinas son nulas o negativas
	if updateData.Tips < 0 {
		return nil, fmt.Errorf(ErrorInvalidTipAmount)
	}

	// Calcular el subtotal
	subtotal := existingInvoice.SubTotal

	// Verificar si los tips son un porcentaje o un valor nominal
	if updateData.Tips <= 1.0 { // Si es menor o igual a 1, se considera un porcentaje
		// Verificar si el porcentaje excede el 10% del subtotal
		if updateData.Tips > 0.1*subtotal {
			return nil, fmt.Errorf(ErrorTipPercentageExceedsLimit)
		}
		existingInvoice.Tips = updateData.Tips * subtotal // Calcular las propinas como un porcentaje del subtotal
	} else { // Si es mayor que 1, se considera un valor nominal y se suma directamente
		existingInvoice.Tips = updateData.Tips
	}
	existingInvoice.ReCalculateTips()
	// Guardar el Invoice actualizado en la base de datos
	updatedInvoice, err := s.repository.UpdateTip(existingInvoice)
	if err != nil {
		return nil, err
	}

	return updatedInvoice, nil
}
