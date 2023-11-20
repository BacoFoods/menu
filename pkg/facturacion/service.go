package facturacion

import (
	"errors"
	"fmt"

	"github.com/BacoFoods/menu/pkg/invoice"
)

const (
	DocumentTypePOS            = "POS"
	DocumentTypeFEIdentified   = "FEIdentified"
	DocumentTypeFEUnidentified = "FEUnidentified"
)

var (
	ErrStoreWithoutConfig = errors.New("store without facturacion config")

	defaultClient = map[string]any{
		"nombre":             "Cosumidor Final",
		"tipoIdentificacion": "nit",
		"identificacion":     "222222222",
		"correo":             "",
	}
)

type FacturacionService struct {
	repository Repository
}

func NewService(repository Repository) *FacturacionService {
	return &FacturacionService{
		repository: repository,
	}
}

func (s *FacturacionService) UpdateConfig(config *FacturacionConfig) (*FacturacionConfig, error) {
	return s.repository.Update(config)
}

func (s *FacturacionService) CreateConfig(config *FacturacionConfig) (*FacturacionConfig, error) {
	if config.DocumentType != DocumentTypePOS && config.DocumentType != DocumentTypeFEIdentified && config.DocumentType != DocumentTypeFEUnidentified {
		return nil, fmt.Errorf("invalid document type")
	}

	if err := s.repository.Create(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *FacturacionService) FindConfig(storeID uint) ([]FacturacionConfig, error) {
	return s.repository.FindByStore(storeID)
}

func (s *FacturacionService) Generate(invoice *invoice.Invoice, docType string, data any) (*invoice.Document, error) {
	switch docType {
	case DocumentTypePOS:
		return s.generatePOS(invoice, data)
	case DocumentTypeFEIdentified:
		return s.generateFEIdentified(invoice, data)
	case DocumentTypeFEUnidentified:
		return s.generateFEUnidentified(invoice)
	default:
		return nil, fmt.Errorf("invalid document type")
	}
}

func (s *FacturacionService) generatePOS(inv *invoice.Invoice, data any) (*invoice.Document, error) {
	if inv.StoreID == nil {
		return nil, errors.New("invoice store id is required")
	}

	config, err := s.repository.FindByStoreAndType(*inv.StoreID, DocumentTypePOS)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, ErrStoreWithoutConfig
	}

	// TODO: validate data
	curNumber := config.LastNumber + 1

	resolution := config.Resolution
	resolution["prefix"] = config.Prefix

	return &invoice.Document{
		DocumentType: DocumentTypePOS,
		Code:         fmt.Sprintf("%s-%d", config.Prefix, curNumber),
		// Client:       data, // TODO:
		Resolution: resolution,
		Seller:     config.Seller,
	}, nil
}

func (s *FacturacionService) generateFEIdentified(invoice *invoice.Invoice, data any) (*invoice.Document, error) {
	return nil, errors.New("TODO: unimplemented")
}

func (s *FacturacionService) generateFEUnidentified(invoice *invoice.Invoice) (*invoice.Document, error) {
	return s.generateFEIdentified(invoice, defaultClient)
}
