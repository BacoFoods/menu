package cashaudit

import (
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/order"
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
	TotalSell         float64    `json:"total_sell"`
	Orders            uint       `json:"orders"`
	Eaters            uint       `json:"eaters"`
	Discounts         float64    `json:"discounts"`
	Surcharges        float64    `json:"surcharges"`
	Tips              float64    `json:"tips"`
	CashIncomes       float64    `json:"cash_incomes"`
	OnlineIncomes     float64    `json:"online_incomes"`
	CardIncomes       float64    `json:"card_incomes"`
}

func GetInvoices(orders []order.Order) []invoice.Invoice {
	invoices := make([]invoice.Invoice, 0)
	for _, order := range orders {
		invoices = append(invoices, order.Invoices...)
	}
	return invoices
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
