package invoice

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
	Cashier string `json:"cashier"`

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
