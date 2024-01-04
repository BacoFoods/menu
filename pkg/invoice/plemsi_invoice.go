package invoice

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/plemsi"
	"github.com/BacoFoods/menu/pkg/shared"
)

const (
	LogPlemsiInvoice = "pkg/invoice/plemsi_invoice"
)

func (i *Invoice) ToPlemsiInvoice(finalCustomer bool) (*plemsi.Invoice, error) {

	plemsiInvoice := plemsi.NewBuilderEndConsumerInvoice()

	// Setting date
	plemsiInvoice.SetDate(i.CreatedAt) // TODO: change to string

	// Setting time
	plemsiInvoice.SetTime(i.CreatedAt) // TODO: change to string

	// Setting prefix
	plemsiInvoice.SetPrefix("SETT") // TODO: define preffix

	// Setting number
	plemsiInvoice.SetNumber(int(i.ID)) //TODO: consecutive invoice ID is valid?

	// Setting order reference
	orderReference, err := plemsi.NewBuilderOrderReference().SetIdOrder(fmt.Sprintf("%v", *i.OrderID)).Build()
	if err != nil {
		shared.LogError("error building plemsi invoice order reference", LogPlemsiInvoice, "ToPlemsiInvoice", err, *i)
		return nil, err
	}

	plemsiInvoice.SetOrderReference(orderReference)

	// Setting account email false, because final customer invoice
	plemsiInvoice.SetSendEmail(true)

	// Setting final customer
	plemsiInvoice.SetIsFinalCustomer(finalCustomer)

	if !finalCustomer {
		// Setting customer
		customer, err := plemsi.NewBuilderCustomer().
			SetName(i.Client.Name).
			SetIdentificationNumber(i.Client.Document).
			SetTypeDocumentIdentificationId(3). // 3 for Cédula TODO: improve to get id from document type
			Build()
		if err != nil {
			shared.LogError("error building plemsi invoice customer", LogPlemsiInvoice, "ToPlemsiInvoice", err, *i)
			return nil, err
		}

		plemsiInvoice.SetCustomer(customer)
	}

	// Setting payment
	if len(i.Payments) == 0 {
		shared.LogError("error invoice without payment", LogPlemsiInvoice, "ToPlemsiInvoice", err, *i)
		return nil, fmt.Errorf(ErrorPlemsiAdapterInvoiceWithoutPayment)
	}

	// TODO: improve to multiples payments
	payment, err := plemsi.NewBuilderPayment().
		SetPaymentFormId(1).
		SetPaymentMethodId(10). // TODO: get id methods to make invoice
		SetPaymentDueDate(i.CreatedAt).
		Build()
	if err != nil {
		shared.LogError("error building plemsi invoice payment", LogPlemsiInvoice, "ToPlemsiInvoice", err, *i)
		return nil, err
	}

	plemsiInvoice.SetPayment(payment)

	// Setting discounts
	plemsiInvoiceDiscounts := make([]plemsi.Discounts, 0)

	if len(i.Discounts) != 0 {
		description := ""
		percentage := 0.0
		for _, discount := range i.Discounts {
			description += discount.Description + " "
			percentage += discount.Percentage
		}

		plemsiDiscount, err := plemsi.NewBuilderDiscounts().
			SetAmount(i.TotalDiscounts).
			SetBaseAmount(i.SubTotal).
			SetAllowancePercent(percentage).
			SetAllowanceChargeReason(description).
			Build()

		if err != nil {
			shared.LogError("error building plemsi invoice discount", LogPlemsiInvoice, "ToPlemsiInvoice", err, i.Discounts)
			return nil, err
		}

		plemsiInvoiceDiscounts = append(plemsiInvoiceDiscounts, *plemsiDiscount)
	}

	plemsiInvoice.SetGeneralAllowances(plemsiInvoiceDiscounts)

	// Setting items and taxes
	plemsiItems := make([]plemsi.Item, 0)
	plemsiTaxes := make([]plemsi.Tax, 0)

	for _, item := range i.Items {

		// Setting taxes
		plemsiItemTaxes := make([]plemsi.ItemTax, 0)

		plemsiItemTax, err := plemsi.NewBuilderItemTax().
			SetTaxId(item.Tax).                   // TODO: get id tax
			SetPercent(item.TaxPercentage * 100). // Plemsi tax percent is 8, not 0.08 for ico
			SetTaxAmount(item.TaxAmount).
			SetTaxableAmount(item.TaxBase).
			Build()
		if err != nil {
			shared.LogError("error building plemsi invoice item tax", LogPlemsiInvoice, "ToPlemsiInvoice", err, item)
			return nil, err
		}

		plemsiItemTaxes = append(plemsiItemTaxes, *plemsiItemTax)

		// Setting item discounts
		plemsiItemDiscounts := make([]plemsi.ItemDiscount, 0)
		for _, discount := range i.Discounts {
			itemDiscount, err := plemsi.NewBuilderItemDiscount().
				SetChargeIndicator(false).
				SetAllowanceChargeReason(discount.Description).
				SetMultiplierFactorNumeric(1). // See plemsi docs
				SetAmount(item.Price * discount.Percentage / 100).
				SetBaseAmount(item.Price).
				Build()
			if err != nil {
				shared.LogError("error building plemsi invoice item discount", LogPlemsiInvoice, "ToPlemsiInvoice", err, item)
				return nil, err
			}
			plemsiItemDiscounts = append(plemsiItemDiscounts, *itemDiscount)
		}

		// Building item
		plemsiItem, err := plemsi.NewBuilderItem().
			SetLineExtensionAmount(item.Price).
			SetTaxTotals(plemsiItemTaxes).
			SetDescription(fmt.Sprintf("%s - %s", item.Name,t com item.Description)).
			SetNotes(item.Comments).
			SetCode(item.SKU).
			SetPriceAmount(item.Price).
			SetBaseQuantity(1).
			SetInvoicedQuantity(1).
			SetAllowanceCharges(plemsiItemDiscounts).
			SetUnitMeasureId(70).           // 70 is ID for unidad, see plemsi docs
			SetTypeItemIdentificationId(1). // 1 is ID for UNSPC, see plemsi docs
			Build()
		if err != nil {
			shared.LogError("error building plemsi invoice item", LogPlemsiInvoice, "ToPlemsiInvoice", err, item)
			return nil, err
		}

		plemsiTax, err := plemsi.NewBuilderTax().
			SetTaxId(item.Tax).                   // TODO: get id tax
			SetPercent(item.TaxPercentage * 100). // Plemsi tax percent is 8, not 0.08 for ico
			SetTaxAmount(item.TaxAmount).
			SetTaxableAmount(item.TaxBase).
			Build()

		if err != nil {
			shared.LogError("error building plemsi invoice taxes", LogPlemsiInvoice, "ToPlemsiInvoice", err, item)
			return nil, err
		}

		plemsiItems = append(plemsiItems, *plemsiItem)
		plemsiTaxes = append(plemsiTaxes, *plemsiTax)
	}

	plemsiInvoice.SetItems(plemsiItems)

	// Setting resolution
	plemsiInvoice.SetResolution(i.ResolutionNumber)

	// Setting allowance total
	plemsiInvoice.SetAllowanceTotal(i.TotalDiscounts)

	// Setting invoice base total
	plemsiInvoice.SetInvoiceBaseTotal(i.SubTotal)

	// Setting invoice tax exclusive total
	plemsiInvoice.SetInvoiceTaxExclusiveTotal(i.BaseTax)

	// Setting invoice tax inclusive total
	plemsiInvoice.SetInvoiceTaxInclusiveTotal(i.SubTotal + i.Taxes)

	// Setting total to pay
	plemsiInvoice.SetTotalToPay(i.SubTotal + i.Taxes - i.TotalDiscounts)

	// Setting all tax totals
	plemsiInvoice.SetAllTaxTotals(plemsiTaxes)

	// Setting Custom Subtotals
	if i.TipAmount != 0 {
		tips, err := plemsi.NewBuilderTip().
			SetAmount(i.TipAmount).
			SetConcept("Propina").
			Build()

		if err != nil {
			shared.LogError("error building plemsi invoice tip", LogPlemsiInvoice, "ToPlemsiInvoice", err, *i)
			return nil, err
		}

		plemsiInvoice.SetCustomSubtotals([]plemsi.Tip{*tips})
	} else {
		plemsiInvoice.SetCustomSubtotals([]plemsi.Tip{})
	}

	// Setting final total to pay
	plemsiInvoice.SetFinalTotalToPay(i.Total)

	return plemsiInvoice.Build()
}
