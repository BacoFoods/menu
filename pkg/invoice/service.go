package invoice

import (
	"fmt"

	clientPKG "github.com/BacoFoods/menu/pkg/client"
)

const LogService = "pkg/invoice/service"

type Service interface {
	Get(invoiceID string) (*Invoice, error)
	Find(filter map[string]any) ([]Invoice, error)
	UpdateTip(value float64, tipType string, invoiceID string) (*Invoice, error)
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

// UpdateTip update 'tips' field of an Invoice .y verifica si es un valor válido
func (s service) UpdateTip(value float64, tipType string, invoiceID string) (*Invoice, error) {

	existingInvoice, err := s.repository.Get(invoiceID)
	if err != nil {
		return nil, err
	}

	if err := existingInvoice.CalculateTip(value, tipType); err != nil {
		return nil, err
	}

	return s.repository.UpdateTip(existingInvoice)
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
