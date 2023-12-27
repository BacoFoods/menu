package plemsi

import (
	"fmt"
	"strings"
	"time"
)

// Invoice

// Builder for Invoice
type Builder struct {
	Invoice
	Errors []error
}

func NewBuilderEndConsumerInvoice() *Builder {
	return new(Builder)
}

func (ib *Builder) SetDate(time *time.Time) *Builder {
	if time == nil {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiDateEmpty))
	}
	ib.Date = time.Format("2006-01-02")
	return ib
}

func (ib *Builder) SetTime(time *time.Time) *Builder {
	if time == nil {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTimeEmpty))
	}
	ib.Time = time.Format("15:04:05")
	return ib
}

func (ib *Builder) SetPrefix(prefix string) *Builder {
	if prefix == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiPrefixEmpty))
	}
	ib.Prefix = prefix
	return ib
}

func (ib *Builder) SetNumber(number int) *Builder {
	if number == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiNumberEmpty))
	}
	ib.Number = number
	return ib
}

func (ib *Builder) SetOrderReference(orderReference *OrderReference) *Builder {
	if orderReference == nil {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiOrderReferenceEmpty))
	}
	ib.OrderReference = *orderReference
	return ib
}

func (ib *Builder) SetSendEmail(sendEmail bool) *Builder {
	ib.SendEmail = sendEmail
	return ib
}

func (ib *Builder) SetCustomer(customer *Customer) *Builder {
	if customer == nil {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiCustomerEmpty))
	}
	ib.Customer = customer
	return ib
}

func (ib *Builder) SetIsFinalCustomer(isFinalCustomer bool) *Builder {
	ib.IsFinalCustomer = isFinalCustomer
	return ib
}

func (ib *Builder) SetPayment(payment *Payment) *Builder {
	if payment == nil {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiPaymentEmpty))
	}
	ib.Payment = *payment
	return ib
}

func (ib *Builder) SetGeneralAllowances(generalAllowances []Discounts) *Builder {
	ib.GeneralAllowances = generalAllowances
	return ib
}

func (ib *Builder) SetItems(items []Item) *Builder {
	if items == nil {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemsEmpty))
	}
	if len(items) == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemsEmpty))
	}
	ib.Items = items
	return ib
}

func (ib *Builder) SetResolution(resolution string) *Builder {
	if resolution == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiResolutionEmpty))
	}
	ib.Resolution = resolution
	return ib
}

func (ib *Builder) SetResolutionText(resolutionText string) *Builder {
	if resolutionText == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiResolutionTextEmpty))
	}
	ib.ResolutionText = resolutionText
	return ib
}

func (ib *Builder) SetHeadNote(headNote string) *Builder {
	if headNote == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiHeadNoteEmpty))
	}
	ib.HeadNote = headNote
	return ib
}

func (ib *Builder) SetFootNote(footNote string) *Builder {
	if footNote == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiFootNoteEmpty))
	}
	ib.FootNote = footNote
	return ib
}

func (ib *Builder) SetNotes(notes string) *Builder {
	if notes == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiNotesEmpty))
	}
	ib.Notes = notes
	return ib
}

func (ib *Builder) SetAllowanceTotal(allowanceTotal float64) *Builder {
	if allowanceTotal < 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiAllowanceTotalNegative))
	}
	ib.AllowanceTotal = allowanceTotal
	return ib
}

func (ib *Builder) SetInvoiceBaseTotal(invoiceBaseTotal float64) *Builder {
	if invoiceBaseTotal == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiInvoiceBaseTotalEmpty))
	}
	ib.InvoiceBaseTotal = invoiceBaseTotal
	return ib
}

func (ib *Builder) SetInvoiceTaxExclusiveTotal(invoiceTaxExclusiveTotal float64) *Builder {
	if invoiceTaxExclusiveTotal == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiInvoiceTaxExclusiveTotalEmpty))
	}
	ib.InvoiceTaxExclusiveTotal = invoiceTaxExclusiveTotal
	return ib
}

func (ib *Builder) SetInvoiceTaxInclusiveTotal(invoiceTaxInclusiveTotal float64) *Builder {
	if invoiceTaxInclusiveTotal == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiInvoiceTaxInclusiveTotalEmpty))
	}
	ib.InvoiceTaxInclusiveTotal = invoiceTaxInclusiveTotal
	return ib
}

func (ib *Builder) SetTotalToPay(totalToPay float64) *Builder {
	if totalToPay == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTotalToPayEmpty))
	}
	ib.TotalToPay = totalToPay
	return ib
}

func (ib *Builder) SetAllTaxTotals(allTaxTotals []Tax) *Builder {
	if allTaxTotals == nil {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiAllTaxTotalsEmpty))
	}
	ib.AllTaxTotals = allTaxTotals
	return ib
}

func (ib *Builder) SetCustomSubtotals(customSubtotals []Tip) *Builder {
	ib.CustomSubtotals = customSubtotals
	return ib
}

func (ib *Builder) SetFinalTotalToPay(finalTotalToPay float64) *Builder {
	if finalTotalToPay == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiFinalTotalToPayEmpty))
	}
	ib.FinalTotalToPay = finalTotalToPay
	return ib
}

func (ib *Builder) Build() (*Invoice, error) {
	if len(ib.Errors) != 0 {
		return nil, ib.Errors[0]
	}

	return &ib.Invoice, nil
}

// Order Reference

// BuilderOrderReference for build a OrderReference
type BuilderOrderReference struct {
	OrderReference
	Errors []error
}

func NewBuilderOrderReference() *BuilderOrderReference {
	return new(BuilderOrderReference)
}

func (ib *BuilderOrderReference) SetIdOrder(idOrder string) *BuilderOrderReference {
	if idOrder == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiOrderReferenceIdOrderEmpty))
	}
	ib.IdOrder = idOrder
	return ib
}

func (ib *BuilderOrderReference) Build() (*OrderReference, error) {
	if len(ib.Errors) != 0 {
		return nil, ib.Errors[0]
	}

	return &ib.OrderReference, nil
}

// Customer

// BuilderCustomer for build a Customer
type BuilderCustomer struct {
	Customer
	Errors []error
}

func NewBuilderCustomer() *BuilderCustomer {
	return new(BuilderCustomer)
}

func (ib *BuilderCustomer) SetIdentificationNumber(identificationNumber string) *BuilderCustomer {
	if identificationNumber == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiCustomerIdentificationNumberEmpty))
	}
	ib.IdentificationNumber = identificationNumber
	return ib
}

func (ib *BuilderCustomer) SetName(name string) *BuilderCustomer {
	if name == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiCustomerNameEmpty))
	}
	ib.Name = name
	return ib
}

func (ib *BuilderCustomer) SetTypeDocumentIdentificationId(typeDocumentIdentificationId int) *BuilderCustomer {
	if typeDocumentIdentificationId == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiCustomerTypeDocumentIdentificationIdEmpty))
	}
	ib.TypeDocumentIdentificationId = typeDocumentIdentificationId
	return ib
}

func (ib *BuilderCustomer) Build() (*Customer, error) {
	if len(ib.Errors) != 0 {
		return nil, ib.Errors[0]
	}

	return &ib.Customer, nil
}

// Payment

// BuilderPayment for build a Payment
type BuilderPayment struct {
	Payment
	Errors []error
}

func NewBuilderPayment() *BuilderPayment {
	return new(BuilderPayment)
}

func (ib *BuilderPayment) SetPaymentFormId(paymentFormId int) *BuilderPayment {
	if paymentFormId == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiPaymentFormIdEmpty))
	}
	ib.PaymentFormId = paymentFormId
	return ib
}

func (ib *BuilderPayment) SetPaymentMethodId(paymentMethodId int) *BuilderPayment {
	if paymentMethodId == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiPaymentMethodIdEmpty))
	}
	ib.PaymentMethodId = paymentMethodId
	return ib
}

func (ib *BuilderPayment) SetPaymentDueDate(paymentDueDate *time.Time) *BuilderPayment {
	if paymentDueDate == nil {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTimeEmpty))
	}
	ib.PaymentDueDate = paymentDueDate.Format("2006-01-02")
	return ib
}

func (ib *BuilderPayment) SetDurationMeasure(durationMeasure string) *BuilderPayment {
	if durationMeasure == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiPaymentDurationMeasureEmpty))
	}
	ib.DurationMeasure = durationMeasure
	return ib
}

func (ib *BuilderPayment) Build() (*Payment, error) {
	if len(ib.Errors) != 0 {
		return nil, ib.Errors[0]
	}

	return &ib.Payment, nil
}

// Discounts

// BuilderDiscounts for build a Discounts
type BuilderDiscounts struct {
	Discounts
	Errors []error
}

func NewBuilderDiscounts() *BuilderDiscounts {
	return new(BuilderDiscounts)
}

func (ib *BuilderDiscounts) SetAllowanceChargeReason(allowanceChargeReason string) *BuilderDiscounts {
	if allowanceChargeReason == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiDiscountsAllowanceChargeReasonEmpty))
	}
	ib.AllowanceChargeReason = allowanceChargeReason
	return ib
}

func (ib *BuilderDiscounts) SetAllowancePercent(allowancePercent float64) *BuilderDiscounts {
	if allowancePercent == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiDiscountsAllowancePercentEmpty))
	}
	ib.AllowancePercent = allowancePercent
	return ib
}

func (ib *BuilderDiscounts) SetAmount(amount float64) *BuilderDiscounts {
	if amount == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiDiscountsAmountEmpty))
	}
	ib.Amount = amount
	return ib
}

func (ib *BuilderDiscounts) SetBaseAmount(baseAmount float64) *BuilderDiscounts {
	if baseAmount == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiDiscountsBaseAmountEmpty))
	}
	ib.BaseAmount = baseAmount
	return ib
}

func (ib *BuilderDiscounts) Build() (*Discounts, error) {
	if len(ib.Errors) != 0 {
		return nil, ib.Errors[0]
	}

	return &ib.Discounts, nil
}

// Item

// BuilderItem for build an Item
type BuilderItem struct {
	Item
	Errors []error
}

func NewBuilderItem() *BuilderItem {
	return new(BuilderItem)
}

func (ib *BuilderItem) SetUnitMeasureId(unitMeasureId int) *BuilderItem {
	if unitMeasureId == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemUnitMeasureIdEmpty))
	}
	ib.UnitMeasureId = unitMeasureId
	return ib
}

func (ib *BuilderItem) SetLineExtensionAmount(lineExtensionAmount float64) *BuilderItem {
	if lineExtensionAmount == 0.0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemLineExtensionAmountEmpty))
	}
	ib.LineExtensionAmount = lineExtensionAmount
	return ib
}

func (ib *BuilderItem) SetFreeOfChargeIndicator(freeOfChargeIndicator bool) *BuilderItem {
	ib.FreeOfChargeIndicator = freeOfChargeIndicator
	return ib
}

func (ib *BuilderItem) SetAllowanceCharges(allowanceCharges []ItemDiscount) *BuilderItem {
	if allowanceCharges == nil {
		// TODO: Check if this is correct
		// ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemAllowanceChargesEmpty))
		ib.AllowanceCharges = []ItemDiscount{}
		return ib
	}
	if len(allowanceCharges) == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemAllowanceChargesEmpty))
	}
	ib.AllowanceCharges = allowanceCharges
	return ib
}

func (ib *BuilderItem) SetTaxTotals(taxTotals []ItemTax) *BuilderItem {
	if taxTotals == nil {
		// TODO: Check if this is correct
		// ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemTaxTotalsEmpty))
		ib.TaxTotals = []ItemTax{}
		return ib
	}
	if len(taxTotals) == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemTaxTotalsEmpty))
	}
	ib.TaxTotals = taxTotals
	return ib
}

func (ib *BuilderItem) SetDescription(description string) *BuilderItem {
	if strings.Trim(description, " ") == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemDescriptionEmpty))
		return ib
	}
	ib.Description = description
	return ib
}

func (ib *BuilderItem) SetNotes(notes string) *BuilderItem {
	ib.Notes = notes
	return ib
}

func (ib *BuilderItem) SetCode(code string) *BuilderItem {
	// TODO: Check if this is correct because this send SKU to plemsi
	if code == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemCodeEmpty))
	}

	ib.Code = code
	return ib
}

func (ib *BuilderItem) SetTypeItemIdentificationId(typeItemIdentificationId int) *BuilderItem {
	if typeItemIdentificationId == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemTypeItemIdentificationIdEmpty))
	}
	ib.TypeItemIdentificationId = typeItemIdentificationId
	return ib
}

func (ib *BuilderItem) SetPriceAmount(priceAmount float64) *BuilderItem {
	if priceAmount == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemPriceAmountEmpty))
	}
	ib.PriceAmount = priceAmount
	return ib
}

func (ib *BuilderItem) SetBaseQuantity(baseQuantity int) *BuilderItem {
	if baseQuantity == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemBaseQuantityEmpty))
	}
	ib.BaseQuantity = baseQuantity
	return ib
}

func (ib *BuilderItem) SetInvoicedQuantity(invoicedQuantity int) *BuilderItem {
	if invoicedQuantity == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemInvoicedQuantityEmpty))
	}
	ib.InvoicedQuantity = invoicedQuantity
	return ib
}

func (ib *BuilderItem) Build() (*Item, error) {
	if len(ib.Errors) > 0 {
		return nil, ib.Errors[0]
	}
	return &ib.Item, nil
}

// ItemDiscount

// BuilderItemDiscount for build an ItemDiscount
type BuilderItemDiscount struct {
	ItemDiscount
	Errors []error
}

func NewBuilderItemDiscount() *BuilderItemDiscount {
	return new(BuilderItemDiscount)
}

func (ib *BuilderItemDiscount) SetChargeIndicator(chargeIndicator bool) *BuilderItemDiscount {
	ib.ChargeIndicator = chargeIndicator
	return ib
}

func (ib *BuilderItemDiscount) SetAllowanceChargeReason(allowanceChargeReason string) *BuilderItemDiscount {
	if allowanceChargeReason == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemDiscountAllowanceChargeReasonEmpty))
	}
	ib.AllowanceChargeReason = allowanceChargeReason
	return ib
}

func (ib *BuilderItemDiscount) SetMultiplierFactorNumeric(multiplierFactorNumeric int) *BuilderItemDiscount {
	if multiplierFactorNumeric == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemDiscountMultiplierFactorNumericEmpty))
	}
	ib.MultiplierFactorNumeric = multiplierFactorNumeric
	return ib
}

func (ib *BuilderItemDiscount) SetAmount(amount float64) *BuilderItemDiscount {
	if amount == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemDiscountPlemsiAmountEmpty))
	}
	ib.Amount = amount
	return ib
}

func (ib *BuilderItemDiscount) SetBaseAmount(baseAmount float64) *BuilderItemDiscount {
	if baseAmount == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemDiscountPlemsiBaseAmountEmpty))
	}
	ib.BaseAmount = baseAmount
	return ib
}

func (ib *BuilderItemDiscount) Build() (*ItemDiscount, error) {
	if len(ib.Errors) != 0 {
		return nil, ib.Errors[0]
	}

	return &ib.ItemDiscount, nil
}

// ItemTax

// BuilderItemTax for build an ItemTax
type BuilderItemTax struct {
	ItemTax
	Errors []error
}

func NewBuilderItemTax() *BuilderItemTax {
	return new(BuilderItemTax)
}

func (ib *BuilderItemTax) SetTaxId(itemTax string) *BuilderItemTax {
	// See plemsi docs codes for taxes
	switch itemTax {
	case "iva":
		ib.TaxId = 1
	case "ico":
		ib.TaxId = 4
	default:
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTaxIdEmpty))
	}
	return ib
}

func (ib *BuilderItemTax) SetPercent(percent float64) *BuilderItemTax {
	if percent == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemTaxPercentEmpty))
	}
	ib.Percent = percent
	return ib
}

func (ib *BuilderItemTax) SetTaxAmount(taxAmount float64) *BuilderItemTax {
	if taxAmount == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemTaxTaxAmountEmpty))
	}
	ib.TaxAmount = taxAmount
	return ib
}

func (ib *BuilderItemTax) SetTaxableAmount(taxableAmount float64) *BuilderItemTax {
	if taxableAmount == 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiItemTaxTaxableAmountEmpty))
	}
	ib.TaxableAmount = taxableAmount
	return ib
}

func (ib *BuilderItemTax) Build() (*ItemTax, error) {
	if len(ib.Errors) != 0 {
		return nil, ib.Errors[0]
	}

	return &ib.ItemTax, nil
}

// Tax

// BuilderTax for build a Tax
type BuilderTax struct {
	Tax
	Errors []error
}

func NewBuilderTax() *BuilderTax {
	return new(BuilderTax)
}

func (ib *BuilderTax) SetTaxId(itemTax string) *BuilderTax {
	// See plemsi docs codes for taxes
	switch itemTax {
	case "iva":
		ib.TaxId = 1
	case "ico":
		ib.TaxId = 4
	default:
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTaxIdEmpty))
	}
	return ib
}

func (ib *BuilderTax) SetTaxAmount(taxAmount float64) *BuilderTax {
	if taxAmount == 0.0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTaxAmountEmpty))
	}
	ib.TaxAmount = taxAmount
	return ib
}

func (ib *BuilderTax) SetPercent(percent float64) *BuilderTax {
	if percent == 0.0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTaxPercentEmpty))
	}
	ib.Percent = percent
	return ib
}

func (ib *BuilderTax) SetTaxableAmount(taxableAmount float64) *BuilderTax {
	if taxableAmount == 0.0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTaxableAmountEmpty))
	}
	ib.TaxableAmount = taxableAmount
	return ib
}

func (ib *BuilderTax) Build() (*Tax, error) {
	if len(ib.Errors) != 0 {
		return nil, ib.Errors[0]
	}

	return &ib.Tax, nil
}

// Tip

// BuilderTip for build a Tip
type BuilderTip struct {
	Tip
	Errors []error
}

func NewBuilderTip() *BuilderTip {
	return new(BuilderTip)
}

func (ib *BuilderTip) SetConcept(concept string) *BuilderTip {
	if concept == "" {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTipConceptEmpty))
	}
	ib.Concept = concept
	return ib
}

func (ib *BuilderTip) SetAmount(amount float64) *BuilderTip {
	if amount < 0 {
		ib.Errors = append(ib.Errors, fmt.Errorf(ErrorPlemsiTipAmountNegative))
	}
	ib.Amount = amount
	return ib
}

func (ib *BuilderTip) Build() (*Tip, error) {
	if len(ib.Errors) != 0 {
		return nil, ib.Errors[0]
	}

	return &ib.Tip, nil
}
