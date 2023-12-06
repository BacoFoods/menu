package invoice

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/plemsi"
	"github.com/BacoFoods/menu/pkg/shared"
)

const (
	LogPlemsiInvoice = "pkg/invoice/plemsi_invoice"
)

func ToPlemsiInvoice(invoice *Invoice) (*plemsi.Invoice, error) {

	plemsiItems := make([]plemsi.Item, 0)

	// Setting items
	for _, item := range invoice.Items {
		plemsiItem, err := plemsi.NewBuilderItem().
			SetDescription(item.Description).
			SetNotes(item.Comments).
			SetCode(item.SKU).
			SetPriceAmount(item.Price).
			SetBaseQuantity(1).
			SetInvoicedQuantity(1).
			Build()

		if err != nil {
			shared.LogError("error building plemsi invoice item", LogPlemsiInvoice, "ToPlemsiInvoice", err, item)
			return nil, err
		}

		plemsiItems = append(plemsiItems, *plemsiItem)
	}

	plemsiInvoice := plemsi.NewBuilderEndConsumerInvoice().SetItems(plemsiItems)

	// Setting date
	plemsiInvoice.SetDate(invoice.CreatedAt) // TODO: change to string
	plemsiInvoice.SetTime(invoice.CreatedAt) // TODO: change to string

	// Setting prefix
	plemsiInvoice.SetPrefix("SETT") // TODO: define preffix

	// Setting number
	plemsiInvoice.SetNumber(int(invoice.ID)) //TODO: consecutive invoice ID is valid?

	// Setting order reference
	orderReference, err := plemsi.NewBuilderOrderReference().SetIdOrder(fmt.Sprintf("%v", *invoice.OrderID)).Build()
	if err != nil {
		shared.LogError("error building plemsi invoice order reference", LogPlemsiInvoice, "ToPlemsiInvoice", err, *invoice)
		return nil, err
	}

	plemsiInvoice.SetOrderReference(orderReference)

	// Setting account email false, because final customer invoice
	plemsiInvoice.SetSendEmail(false)

	// Setting final customer
	plemsiInvoice.SetIsFinalCustomer(true)

	// Setting payment
	if len(invoice.Payments) == 0 {
		shared.LogError("error invoice without payment", LogPlemsiInvoice, "ToPlemsiInvoice", err, *invoice)
		return nil, fmt.Errorf(ErrorPlemsiAdapterInvoiceWithoutPayment)
	}

	// TODO: improve to multiples payments
	invoicePayment := invoice.Payments[0]
	invoicePaymentID := int(invoicePayment.ID)
	payment, err := plemsi.NewBuilderPayment().
		SetPaymentFormId(invoicePaymentID).
		SetPaymentMethodId(1). // TODO: get id methods to make invoice
		SetPaymentDueDate(fmt.Sprint(invoicePaymentID)).
		Build()
	if err != nil {
		shared.LogError("error building plemsi invoice payment", LogPlemsiInvoice, "ToPlemsiInvoice", err, *invoice)
		return nil, err
	}

	plemsiInvoice.SetPayment(payment)

	// Setting discounts
	plemsiInvoiceDiscounts := make([]plemsi.Discounts, 0)

	if len(invoice.Discounts) != 0 {
		for _, discount := range invoice.Discounts {
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

}
