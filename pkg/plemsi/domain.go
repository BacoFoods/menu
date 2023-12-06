package plemsi

const (
	ErrorPlemsiDateEmpty                     = "error plemsi adapter date is empty"
	ErrorPlemsiTimeEmpty                     = "error plemsi adapter time is empty"
	ErrorPlemsiPrefixEmpty                   = "error plemsi adapter prefix is empty"
	ErrorPlemsiNumberEmpty                   = "error plemsi adapter number is empty"
	ErrorPlemsiOrderReferenceEmpty           = "error plemsi adapter order reference is empty"
	ErrorPlemsiCustomerEmpty                 = "error plemsi adapter customer is empty"
	ErrorPlemsiPaymentEmpty                  = "error plemsi adapter payment is empty"
	ErrorPlemsiGeneralAllowancesEmpty        = "error plemsi adapter general allowances is empty"
	ErrorPlemsiItemsEmpty                    = "error plemsi adapter items is empty"
	ErrorPlemsiResolutionEmpty               = "error plemsi adapter resolution is empty"
	ErrorPlemsiResolutionTextEmpty           = "error plemsi adapter resolution text is empty"
	ErrorPlemsiHeadNoteEmpty                 = "error plemsi adapter head note is empty"
	ErrorPlemsiFootNoteEmpty                 = "error plemsi adapter foot note is empty"
	ErrorPlemsiNotesEmpty                    = "error plemsi adapter notes is empty"
	ErrorPlemsiAllowanceTotalEmpty           = "error plemsi adapter allowance total is empty"
	ErrorPlemsiInvoiceBaseTotalEmpty         = "error plemsi adapter invoice base total is empty"
	ErrorPlemsiInvoiceTaxExclusiveTotalEmpty = "error plemsi adapter invoice tax exclusive total is empty"
	ErrorPlemsiInvoiceTaxInclusiveTotalEmpty = "error plemsi adapter invoice tax inclusive total is empty"
	ErrorPlemsiTotalToPayEmpty               = "error plemsi adapter total to pay is empty"
	ErrorPlemsiAllTaxTotalsEmpty             = "error plemsi adapter all tax totals is empty"
	ErrorPlemsiCustomSubtotalsEmpty          = "error plemsi adapter custom subtotals is empty"
	ErrorPlemsiFinalTotalToPayEmpty          = "error plemsi adapter final total to pay is empty"
	ErrorPlemsiOrderReferenceIdOrderEmpty    = "error plemsi adapter order reference id order is empty"

	ErrorPlemsiCustomerIdentificationNumberEmpty         = "error plemsi adapter customer identification number is empty"
	ErrorPlemsiCustomerNameEmpty                         = "error plemsi adapter customer name is empty"
	ErrorPlemsiCustomerTypeDocumentIdentificationIdEmpty = "error plemsi adapter customer type document identification id is empty"

	ErrorPlemsiPaymentFormIdEmpty          = "error plemsi adapter payment form id is empty"
	ErrorPlemsiPaymentMethodIdEmpty        = "error plemsi adapter payment method id is empty"
	ErrorPlemsiPaymentDueDateEmpty         = "error plemsi adapter payment due date is empty"
	ErrorPlemsiPaymentDurationMeasureEmpty = "error plemsi adapter payment duration measure is empty"

	ErrorPlemsiDiscountsAllowanceChargeReasonEmpty = "error plemsi adapter discounts allowance charge reason is empty"
	ErrorPlemsiDiscountsAllowancePercentEmpty      = "error plemsi adapter discounts allowance percent is empty"
	ErrorPlemsiDiscountsAmountEmpty                = "error plemsi adapter discounts amount is empty"
	ErrorPlemsiDiscountsBaseAmountEmpty            = "error plemsi adapter discounts base amount is empty"

	ErrorPlemsiItemUnitMeasureIdEmpty                   = "error plemsi adapter item unit measure id is empty"
	ErrorPlemsiItemLineExtensionAmountEmpty             = "error plemsi adapter item line extension amount is empty"
	ErrorPlemsiItemAllowanceChargesEmpty                = "error plemsi adapter item allowance charges is empty"
	ErrorPlemsiItemTaxTotalsEmpty                       = "error plemsi adapter item tax totals is empty"
	ErrorPlemsiItemDescriptionEmpty                     = "error plemsi adapter item description is empty"
	ErrorPlemsiItemNotesEmpty                           = "error plemsi adapter item notes is empty"
	ErrorPlemsiItemCodeEmpty                            = "error plemsi adapter item code is empty"
	ErrorPlemsiItemTypeItemIdentificationIdEmpty        = "error plemsi adapter item type item identification id is empty"
	ErrorPlemsiItemPriceAmountEmpty                     = "error plemsi adapter item price amount is empty"
	ErrorPlemsiItemBaseQuantityEmpty                    = "error plemsi adapter item base quantity is empty"
	ErrorPlemsiItemInvoicedQuantityEmpty                = "error plemsi adapter item invoiced quantity is empty"
	ErrorPlemsiItemDiscountAllowanceChargeReasonEmpty   = "error plemsi adapter item discount allowance charge reason is empty"
	ErrorPlemsiItemDiscountMultiplierFactorNumericEmpty = "error plemsi adapter item discount multiplier factor numeric is empty"
	ErrorPlemsiItemDiscountPlemsiAmountEmpty            = "error plemsi adapter item discount amount is empty"
	ErrorPlemsiItemDiscountPlemsiBaseAmountEmpty        = "error plemsi adapter item discount base amount is empty"
	ErrorPlemsiItemTaxTaxIdEmpty                        = "error plemsi adapter item tax tax id is empty"
	ErrorPlemsiItemTaxPercentEmpty                      = "error plemsi adapter item tax percent is empty"
	ErrorPlemsiItemTaxTaxAmountEmpty                    = "error plemsi adapter item tax tax amount is empty"
	ErrorPlemsiItemTaxTaxableAmountEmpty                = "error plemsi adapter item tax taxable amount is empty"

	ErrorPlemsiTaxIdEmpty         = "error plemsi adapter tax tax id is empty"
	ErrorPlemsiTaxAmountEmpty     = "error plemsi adapter tax tax amount is empty"
	ErrorPlemsiTaxPercentEmpty    = "error plemsi adapter tax percent is empty"
	ErrorPlemsiTaxableAmountEmpty = "error plemsi adapter taxable amount is empty"

	ErrorPlemsiTipConceptEmpty = "error plemsi adapter tip concept is empty"
	ErrorPlemsiTipAmountEmpty  = "error plemsi adapter tip amount is empty"

	ErrorPlemsiEndConsumerInvoice = "error plemsi adapter end consumer invoice integration"
	ErrorPlemsiTestConnection     = "error plemsi adapter test connection"
	ErrorPlemsiEmptyInvoice       = "error plemsi empty invoice"
)

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
