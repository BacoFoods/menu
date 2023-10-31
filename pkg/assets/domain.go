package assets

import (
	"time"
)

type Asset struct {
	ID                     uint       `json:"id" swaggertype:"integer" example:"1" swaggerignore:"true"`
	Reference              string     `json:"reference" swaggertype:"string" example:"REF001"`
	Placa                  string     `json:"placa" swaggertype:"string" example:"PLACA123"`
	AssetName              *string    `json:"asset_name" swaggertype:"string" example:"Asset Name"`
	AssetNameSiesa         *string    `json:"asset_name_siesa" swaggertype:"string" example:"Asset Siesa"`
	OperationCodeSiesa     *string    `json:"operation_code_siesa" swaggertype:"string" example:"OP123"`
	OperationNameSiesa     *string    `json:"operation_name_siesa" swaggertype:"string" example:"Operation"`
	PurchaseDate           *time.Time `json:"purchase_date" swaggertype:"string" example:"2022-10-10"`
	PurchaseInvoice        *string    `json:"purchase_invoice" swaggertype:"string" example:"INV001"`
	ProviderNIT            *string    `json:"provider_nit" swaggertype:"string" example:"123456789"`
	ProviderName           *string    `json:"provider_name" swaggertype:"string" example:"Provider Inc."`
	OfficialPlaca          string     `json:"official_placa" swaggertype:"string" example:"OFFPLACA123"`
	CurrentLocation        string     `json:"current_location" swaggertype:"string" example:"Location A"`
	CurrentCostCenter      int        `json:"current_cost_center" swaggertype:"integer" example:"1001"`
	Custodian              string     `json:"custodian" swaggertype:"string" example:"John Doe"`
	Invoice                string     `json:"invoice" swaggertype:"string" example:"INV123"`
	InvoiceLink            string     `json:"invoice_link" swaggertype:"string" example:"http://example.com/invoice/INV123"`
	ActualPurchaseDate     time.Time  `json:"actual_purchase_date" swaggertype:"string" example:"2022-10-01"`
	Contract               string     `json:"contract" swaggertype:"string" example:"CON001"`
	Category               string     `json:"category" swaggertype:"string" example:"Category A"`
	Renting                bool       `json:"renting" swaggertype:"boolean" example:"false"`
	Price                  float64    `json:"price" swaggertype:"number" example:"1000.50"`
	Barcode                string     `json:"barcode" swaggertype:"string" example:"BARCODE001"`
	Type                   string     `json:"type" swaggertype:"string" example:"SIESA"`
	InvoiceInitialLocation *string    `json:"invoice_initial_location" swaggertype:"string" example:"Location A"`
	CreatedAt              time.Time  `json:"created_at" swaggerignore:"true"`
	UpdatedAt              time.Time  `json:"updated_at" swaggerignore:"true"`
}

func (Asset) TableName() string {
	return "assets"
}
