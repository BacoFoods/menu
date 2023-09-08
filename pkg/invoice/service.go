package invoice

import "fmt"

const LogService = "pkg/invoice/service"

type Service interface {
	Get(invoiceID string) (*Invoice, error)
	Update(invoice *Invoice, discounts []uint, Surcharges []uint) (*Invoice, error)
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

	fmt.Println("existingInvoice",existingInvoice)
	fmt.Println("existingInvoice",existingInvoice)
	fmt.Println("discounts",discounts)
	fmt.Println("DiscountsInvoice existente",existingInvoice.Discounts)

	if updateData.Tips > 0.1*existingInvoice.SubTotal {
		return nil, fmt.Errorf("tips cannot be greater than 10 percent of subtotal")
	}

	fmt.Println("updateData",updateData.Discounts)

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

		fmt.Println("surcharge completo",surcharge)
		
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
