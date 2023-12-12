package plemsi

type Invoice struct {
	Date                     string         `json:"date"`
	Time                     string         `json:"time"`
	Prefix                   string         `json:"prefix"`
	Number                   int            `json:"number"`
	OrderReference           OrderReference `json:"orderReference"`
	SendEmail                bool           `json:"send_email"`
	Customer                 Customer       `json:"customer"`
	IsFinalCustomer          bool           `json:"isFinalCustomer"`
	Payment                  Payment        `json:"payment"`
	GeneralAllowances        []Discounts    `json:"generalAllowances"`
	Items                    []Item         `json:"items"`
	Resolution               string         `json:"resolution"`
	ResolutionText           string         `json:"resolutionText"`
	HeadNote                 string         `json:"head_note"`
	FootNote                 string         `json:"foot_note"`
	Notes                    string         `json:"notes"`
	AllowanceTotal           int            `json:"allowanceTotal"`
	InvoiceBaseTotal         int            `json:"invoiceBaseTotal"`
	InvoiceTaxExclusiveTotal int            `json:"invoiceTaxExclusiveTotal"`
	InvoiceTaxInclusiveTotal int            `json:"invoiceTaxInclusiveTotal"`
	TotalToPay               int            `json:"totalToPay"`
	AllTaxTotals             []Tax          `json:"allTaxTotals"`
	CustomSubtotals          []Tip          `json:"customSubtotals"`
	FinalTotalToPay          int            `json:"finalTotalToPay"`
}

type Customer struct {
	IdentificationNumber         string `json:"identification_number"`
	Name                         string `json:"name"`
	TypeDocumentIdentificationId int    `json:"type_document_identification_id"`
}

type OrderReference struct {
	IdOrder string `json:"id_order"`
}

type Payment struct {
	PaymentFormId   int    `json:"payment_form_id"`
	PaymentMethodId int    `json:"payment_method_id"`
	PaymentDueDate  string `json:"payment_due_date"`
	DurationMeasure string `json:"duration_measure"`
}

type Discounts struct {
	AllowanceChargeReason string  `json:"allowance_charge_reason"`
	AllowancePercent      float64 `json:"allowance_percent"`
	Amount                float64 `json:"amount"`
	BaseAmount            float64 `json:"base_amount"`
}

type Item struct {
	UnitMeasureId            int            `json:"unit_measure_id"`
	LineExtensionAmount      int            `json:"line_extension_amount"`
	FreeOfChargeIndicator    bool           `json:"free_of_charge_indicator"`
	AllowanceCharges         []ItemDiscount `json:"allowance_charges"`
	TaxTotals                []ItemTax      `json:"tax_totals"`
	Description              string         `json:"description"`
	Notes                    string         `json:"notes"`
	Code                     string         `json:"code"`
	TypeItemIdentificationId int            `json:"type_item_identification_id"`
	PriceAmount              float64        `json:"price_amount"`
	BaseQuantity             int            `json:"base_quantity"`
	InvoicedQuantity         int            `json:"invoiced_quantity"`
}

type ItemDiscount struct {
	ChargeIndicator         bool   `json:"charge_indicator"`
	AllowanceChargeReason   string `json:"allowance_charge_reason"`
	MultiplierFactorNumeric int    `json:"multiplier_factor_numeric"`
	Amount                  int    `json:"amount"`
	BaseAmount              int    `json:"base_amount"`
}

type ItemTax struct {
	TaxId         int `json:"tax_id"`
	Percent       int `json:"percent"`
	TaxAmount     int `json:"tax_amount"`
	TaxableAmount int `json:"taxable_amount"`
}

type Tax struct {
	TaxId         int `json:"tax_id"`
	TaxAmount     int `json:"tax_amount"`
	Percent       int `json:"percent"`
	TaxableAmount int `json:"taxable_amount"`
}

type Tip struct {
	Concept string `json:"concept"`
	Amount  int    `json:"amount"`
}
