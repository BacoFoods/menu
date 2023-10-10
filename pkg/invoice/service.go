package invoice

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"

	clientPKG "github.com/BacoFoods/menu/pkg/client"
)

const LogService = "pkg/invoice/service"

type Service interface {
	Get(invoiceID string) (*Invoice, error)
	Find(filter map[string]any) ([]Invoice, error)
	UpdateTip(value float64, tipType string, invoiceID string) (*Invoice, error)
	AddClient(invoiceID string, clientID string) (*Invoice, error)
	RemoveClient(invoiceID string, clientID string) (*Invoice, error)
	Separate(invoiceID string, invoices [][]uint) ([]Invoice, error)
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

// Separate separates an invoice into multiple invoices.
func (s service) Separate(invoiceID string, invoices [][]uint) ([]Invoice, error) {
	invoiceDB, err := s.repository.Get(invoiceID)
	if err != nil {
		return nil, err
	}

	mapItems := invoiceDB.MapItems()

	newInvoices := make([]Invoice, 0)
	for _, invoice := range invoices {
		newInvoice := Invoice{
			OrderID:         invoiceDB.OrderID,
			BrandID:         invoiceDB.BrandID,
			StoreID:         invoiceDB.StoreID,
			ChannelID:       invoiceDB.ChannelID,
			TableID:         invoiceDB.TableID,
			Items:           nil,
			Discounts:       invoiceDB.Discounts,
			Surcharges:      invoiceDB.Surcharges,
			SubTotal:        0,
			TotalDiscounts:  0,
			TotalSurcharges: 0,
			Tips:            0,
			BaseTax:         0,
			Taxes:           0,
			Total:           0,
			Payments:        nil,
			Client:          clientPKG.DefaultClient(),
		}
		for _, itemID := range invoice {
			item, ok := mapItems[itemID]
			if !ok {
				err := fmt.Errorf(ErrorItemNotFound)
				shared.LogError("error separating invoice", LogService, "Separate", err, itemID)
				return nil, err
			}
			delete(mapItems, itemID) // remove item from map to validate that all items are separated
			newInvoice.Items = append(newInvoice.Items, item)
			newInvoice.SubTotal += item.Price
			tax := newInvoice.SubTotal * TaxPercentage
			baseTax := newInvoice.SubTotal - tax
			newInvoice.BaseTax = baseTax
			newInvoice.Total = newInvoice.SubTotal + tax
		}
		newInvoices = append(newInvoices, newInvoice)
	}

	if len(mapItems) != 0 {
		err := fmt.Errorf(ErrorInvoiceSeparatingNotEnoughItems)
		shared.LogError("error separating invoice", LogService, "Separate", err, mapItems)
		return nil, err
	}

	invoicesBatch, err := s.repository.CreateBatch(newInvoices)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Delete(invoiceID); err != nil {
		return nil, err
	}

	return invoicesBatch, nil
}
