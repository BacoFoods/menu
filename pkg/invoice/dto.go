package invoice

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"time"
)

const LogDTO = "pkg/invoice/dto"

type DTOPrintable struct {
	StoreName    string `json:"store_name"`
	StoreAddress string `json:"store_address"`
	StorePhone   string `json:"store_phone"`

	BrandName     string `json:"brand_name"`
	BrandDocument string `json:"brand_document"`
	BrandCity     string `json:"brand_city"`
	//-------------------------------------------//
	// TODO: Pending -> Tipo de factura + Prefijo de factura + Código de Factura + Incremental => FACTURA ELECTRÓNICA DE VENTA No. FE90 - 01
	Date    string `json:"date"`
	Waiter  string `json:"waiter"`
	Cashier string `json:"shift"`

	ClientName     string `json:"client_name"`
	ClientDocument string `json:"client_document"`
	ClientEmail    string `json:"client_email"`
	ClientAddress  string `json:"client_address"`
	//-------------------------------------------//
	// TODO: Pending -> Identificador de Tipo de Pedido + Identificador de usuario
	OrderID string `json:"order_id"`
	Table   string `json:"table_name"`
	//-------------------------------------------//
	Items []DTOPrintableItem `json:"items" gorm:"-"`
	//-------------------------------------------//
	Subtotal  float64 `json:"subtotal"`
	Discount  float64 `json:"discount"`
	Tip       string  `json:"tip"`
	TipAmount float64 `json:"tip_amount"`
	Surcharge float64 `json:"surcharge"`
	Total     float64 `json:"total"`
	//-------------------------------------------//
	BaseTax  float64 `json:"base_tax"`
	TaxIPOCO float64 `json:"tax_ipoco"`
	TaxIVA   float64 `json:"tax_iva"`
	//-------------------------------------------//
}

type DTOPrintableItem struct {
	Name     string  `json:"name"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
	Total    float64 `json:"total"`
}

type DTOResolution struct {
	BrandID        *uint  `json:"brand_id"`
	StoreID        *uint  `json:"store_id" validate:"required"`
	DateFrom       string `json:"date_from" example:"2021-01-01" format:"2006-01-02"`
	DateTo         string `json:"date_to" example:"2021-01-01" format:"2006-01-02"`
	Prefix         string `json:"prefix"`
	From           *int   `json:"from" validate:"required"`
	To             *int   `json:"to" validate:"required"`
	Current        *int   `json:"current"`
	Resolution     string `json:"resolution" validate:"required"`
	ResolutionDate string `json:"resolution_date" example:"2021-01-01" format:"2006-01-02"`
	TypeResolution string `json:"type_resolution" example:"electrónica" enum:"electrónica, papel, talonario" validate:"required"`
	TypeDocument   string `json:"type_document" example:"factura, nota crédito, nota débito" enum:"factura, nota crédito, nota débito" validate:"required"`
}

func (dto *DTOResolution) ToResolution() (*Resolution, error) {
	dateFrom, err := time.Parse("2006-01-02", dto.DateFrom)
	if err != nil {
		shared.LogError("error parsing date", LogDTO, "ToResolution", err, dto.DateFrom)
		return nil, err
	}

	dateTo, err := time.Parse("2006-01-02", dto.DateTo)
	if err != nil {
		shared.LogError("error parsing date", LogDTO, "ToResolution", err, dto.DateTo)
		return nil, err
	}

	resolutionDate, err := time.Parse("2006-01-02", dto.ResolutionDate)
	if err != nil {
		shared.LogError("error parsing date", LogDTO, "ToResolution", err, dto.ResolutionDate)
		return nil, err
	}

	return &Resolution{
		BrandID:        dto.BrandID,
		StoreID:        dto.StoreID,
		DateFrom:       &dateFrom,
		DateTo:         &dateTo,
		Prefix:         dto.Prefix,
		From:           dto.From,
		To:             dto.To,
		Current:        dto.Current,
		Resolution:     dto.Resolution,
		ResolutionDate: &resolutionDate,
		TypeResolution: dto.TypeResolution,
		TypeDocument:   dto.TypeDocument,
	}, nil
}
