package cashaudit

import (
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/payment"
	"time"
)

const (
	ErrorCashAuditGetting          = "error getting cash audit"
	ErrorCashAuditGettingStore     = "error getting store"
	ErrorCashAuditGettingStoreID   = "error getting store id"
	ErrorCashAuditGettingLastShift = "error getting last shift"
	ErrorCashAuditGettingOrders    = "error getting orders"
)

type CashAudit struct {
	StoreName         string     `json:"store_name"`
	ShiftOpen         *time.Time `json:"shift_open"`
	ShiftStartBalance float64    `json:"shift_start_balance"`
	ShiftClose        *time.Time `json:"shift_close"`
	ShiftEndBalance   float64    `json:"shift_end_balance"`
	Orders            uint       `json:"orders"`
	Eaters            uint       `json:"eaters"`
	Discounts         float64    `json:"discounts"`
	Surcharges        float64    `json:"surcharges"`
	Tips              float64    `json:"tips"`
	// Calculated section is for the calculated values from invoices
	TotalSellCalculated     float64 `json:"total_sell"`
	CashIncomesCalculated   float64 `json:"cash_incomes"`
	OnlineIncomesCalculated float64 `json:"online_incomes"`
	CardIncomesCalculated   float64 `json:"card_incomes"`
	// Reported section is for the reported values from cashier
	TotalSellReported     float64 `json:"total_sell_reported"`
	CashIncomesReported   float64 `json:"cash_incomes_reported"`
	OnlineIncomesReported float64 `json:"online_incomes_reported"`
	CardIncomesReported   float64 `json:"card_incomes_reported"`
}

func GetInvoices(orders []order.Order) []invoice.Invoice {
	invoices := make([]invoice.Invoice, 0)
	for _, order := range orders {
		invoices = append(invoices, order.Invoices...)
	}
	return invoices
}

func GetPayments(invoices []invoice.Invoice) []payment.Payment {
	payments := make([]payment.Payment, 0)
	for _, invoice := range invoices {
		payments = append(payments, invoice.Payments...)
	}
	return payments
}

func GetTotalEaters(orders []order.Order) uint {
	eaters := 0
	for _, order := range orders {
		eaters += order.Seats
	}
	return uint(eaters)
}

func GetTotalSell(invoices []invoice.Invoice) float64 {
	total := 0.0
	for _, invoice := range invoices {
		total += invoice.Total
	}

	return total
}

func GetTotalCashIncomes(payments []payment.Payment) float64 {
	total := 0.0
	for _, paymnt := range payments {
		if paymnt.Method == "cash" {
			total += float64(paymnt.TotalValue)
		}
	}

	return total
}

func GetTotalOnlineIncomes(payments []payment.Payment) float64 {
	total := 0.0
	for _, paymnt := range payments {
		if paymnt.Method == "yuno" || paymnt.Method == "paylot" {
			total += float64(paymnt.TotalValue)
		}
	}

	return total
}

func GetTotalCardIncomes(payments []payment.Payment) float64 {
	total := 0.0
	for _, paymnt := range payments {
		if paymnt.Method == "card" {
			total += float64(paymnt.TotalValue)
		}
	}

	return total
}

func GetTotalDiscounts(invoices []invoice.Invoice) float64 {
	total := 0.0
	for _, invoice := range invoices {
		total += invoice.TotalDiscounts
	}

	return total
}

func GetTotalSourcharges(invoices []invoice.Invoice) float64 {
	total := 0.0
	for _, invoice := range invoices {
		total += invoice.TotalSurcharges
	}

	return total
}

func GetTotalTips(invoices []invoice.Invoice) float64 {
	total := 0.0
	for _, invoice := range invoices {
		total += invoice.TipAmount
	}

	return total
}
