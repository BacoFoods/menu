package invoice

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/plemsi"
	"github.com/BacoFoods/menu/pkg/shared"
)

const (
	LogPlemsiInvoice = "pkg/invoice/plemsi_invoice"
)

func (i *Invoice) ToPlemsiInvoice(resolution string) (*plemsi.Invoice, error) {

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
	plemsiInvoice.SetSendEmail(false)

	// Setting final customer
	plemsiInvoice.SetIsFinalCustomer(true)

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
		for _, discount := range i.Discounts {
			plemsiDiscount, err := plemsi.NewBuilderDiscounts().
				SetAmount(discount.Amount).
				SetBaseAmount(discount.Percentage).
				SetAllowancePercent(discount.Percentage).
				SetAllowanceChargeReason(discount.Description).
				Build()

			if err != nil {
				shared.LogError("error building plemsi invoice discount", LogPlemsiInvoice, "ToPlemsiInvoice", err, discount)
				return nil, err
			}

			plemsiInvoiceDiscounts = append(plemsiInvoiceDiscounts, *plemsiDiscount)
		}
	}

	plemsiInvoice.SetGeneralAllowances(plemsiInvoiceDiscounts)

	// Setting items
	plemsiItems := make([]plemsi.Item, 0)

	for _, item := range i.Items {
		plemsiItem, err := plemsi.NewBuilderItem().
			SetTaxTotals(nil).
			SetDescription(item.Description).
			SetNotes(item.Comments).
			SetCode(item.SKU).
			SetPriceAmount(item.Price).
			SetBaseQuantity(1).
			SetInvoicedQuantity(1).
			SetAllowanceCharges(nil).
			SetUnitMeasureId(1).                               // TODO: get id unit measure
			SetTypeItemIdentificationId(int(*item.ProductID)). // TODO: get id type item identification
			Build()

		if err != nil {
			shared.LogError("error building plemsi invoice item", LogPlemsiInvoice, "ToPlemsiInvoice", err, item)
			return nil, err
		}

		plemsiItems = append(plemsiItems, *plemsiItem)
	}

	plemsiInvoice.SetItems(plemsiItems)

	// Setting resolution
	plemsiInvoice.SetResolution(resolution)

	// Setting allowance total
	plemsiInvoice.SetAllowanceTotal(int(i.TotalDiscounts))

	// Setting invoice base total
	plemsiInvoice.SetInvoiceBaseTotal(int(i.SubTotal))

	// Setting invoice tax exclusive total
	plemsiInvoice.SetInvoiceTaxExclusiveTotal(int(i.SubTotal))

	// Setting invoice tax inclusive total
	plemsiInvoice.SetInvoiceTaxInclusiveTotal(int(i.BaseTax + i.Taxes))

	// Setting total to pay
	plemsiInvoice.SetTotalToPay(int(i.Total))

	// Setting all tax totals
	plemsiInvoice.SetAllTaxTotals(nil)

	// Setting Custom Subtotals
	if i.TipAmount != 0 {
		tips, err := plemsi.NewBuilderTip().
			SetAmount(int(i.TipAmount)).
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
	plemsiInvoice.SetFinalTotalToPay(int(i.Total))

	return plemsiInvoice.Build()
}
