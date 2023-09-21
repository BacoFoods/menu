package invoice

import (
	"fmt"

	clientPKG "github.com/BacoFoods/menu/pkg/client"
)

const LogService = "pkg/invoice/service"

type Service interface {
	Get(invoiceID string) (*Invoice, error)
	Find(filter map[string]any) ([]Invoice, error)
	UpdateTip(invoice *Invoice) (*Invoice, error)
	AddClient(invoiceID string, clientID string) (*Invoice, error)
	RemoveClient(invoiceID string, clientID string) (*Invoice, error)
}

type service struct {
	repository       Repository
	clientRepository clientPKG.Repository
}

func NewService(repository Repository, clientRepository clientPKG.Repository) service {
	return service{repository, clientRepository}
}

// Get returns a single Invoice object by ID.
func (s service) Get(invoiceID string) (*Invoice, error) {
	return s.repository.Get(invoiceID)
}

// Find returns a list of Invoice objects.
func (s service) Find(filter map[string]any) ([]Invoice, error) {
	return s.repository.Find(filter)
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

// AddClient adds a client to an invoice.
func (s service) AddClient(invoiceID string, clientID string) (*Invoice, error) {
	invoice, err := s.repository.Get(invoiceID)
	if err != nil {
		return nil, err
	}

	client, err := s.clientRepository.Get(clientID)
	if err != nil {
		return nil, err
	}

	invoice.ClientID = &client.ID
	return s.repository.CreateUpdate(invoice)
}

// RemoveClient removes a client from an invoice.
func (s service) RemoveClient(invoiceID string, clientID string) (*Invoice, error) {
	invoice, err := s.repository.Get(invoiceID)
	if err != nil {
		return nil, err
	}

	client, err := s.clientRepository.Get(clientID)
	if err != nil {
		return nil, err
	}

	if invoice.ClientID == nil || *invoice.ClientID != client.ID {
		return nil, fmt.Errorf(ErrorInvoiceWrongClient)
	}

	invoice.ClientID = nil
	return s.repository.CreateUpdate(invoice)
}
