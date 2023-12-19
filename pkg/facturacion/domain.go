package facturacion

import (
	"github.com/BacoFoods/menu/internal"
)

type FacturacionConfig struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	StoreID      uint   `json:"store_id" gorm:"uniqueIndex:idx_store_document_type"`
	DocumentType string `json:"document_type" gorm:"uniqueIndex:idx_store_document_type"`
	Prefix       string `json:"prefix"`

	// format should be  { "name": "string", "nit": "string", "address": "string", "city": "string", "regimen": "string", "phone": "string" }
	Seller internal.JSONMap `json:"seller" gorm:"type:jsonb"`

	// format should be  { "from": <num>, "to": <num>,  "number": "string", "date_init": "DD-MM-YYYY", "date_end": "DD-MM-YYYY" }
	Resolution internal.JSONMap `json:"resolution" gorm:"type:jsonb"`

	// starts at 0
	LastNumber uint `json:"last_number"`
}

type FacturacionConfigRepository interface {
	Update(config *FacturacionConfig) (*FacturacionConfig, error)
	Create(config *FacturacionConfig) error
	FindByStoreAndType(storeID uint, docType string) (*FacturacionConfig, error)
	FindByStoreAndTypeAndIncrement(storeID uint, docType string) (*FacturacionConfig, error)
	FindByStore(storeID uint) ([]FacturacionConfig, error)
}
