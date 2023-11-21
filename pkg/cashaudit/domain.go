package cashaudit

import (
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/payment"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorCashAuditGetting          = "error getting cash audit"
	ErrorCashAuditGettingStore     = "error getting store"
	ErrorCashAuditGettingStoreID   = "error getting store id"
	ErrorCashAuditGettingLastShift = "error getting last shift"
	ErrorCashAuditGettingOrders    = "error getting orders"
	ErrorCashAuditNotOrdersFound   = "error not orders found"
)

type TotalIncome struct {
	ID          uint           `json:"id"`
	CashAuditID *uint          `json:"cash_audit_id"`
	Income      float64        `json:"income" gorm:"precision:18;scale:4"`
	Origin      string         `json:"origin"`
	Type        string         `json:"type" enums:"cash,online,card,other"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type CashAudit struct {
	ID                uint       `json:"id"`
	StoreName         string     `json:"store_name"`
	ShiftOpen         *time.Time `json:"shift_open"`
	ShiftStartBalance float64    `json:"shift_start_balance"`
	ShiftClose        *time.Time `json:"shift_close"`
	ShiftEndBalance   float64    `json:"shift_end_balance"`
	Orders            uint       `json:"orders"`
	Eaters            uint       `json:"eaters"`
	Discounts         float64    `json:"discounts"`
	Surcharges        float64    `json:"surcharges"`
	TotalSell         float64    `json:"total_sell"`
	BruteSell         float64    `json:"brute_sell"`
	// Calculated section is for the calculated values from invoices
	//TotalSellCalculated     float64 `json:"total_sell"`
	//CashIncomesCalculated   float64 `json:"cash_incomes"`
	//OnlineIncomesCalculated float64 `json:"online_incomes"`
	//CardIncomesCalculated   float64 `json:"card_incomes"`
	//OtherIncomeCalculated   float64 `json:"other_income"`
	TotalIncomes []TotalIncome `json:"total_incomes" gorm:"foreignKey:CashAuditID"`
	// Reported section is for the reported values from cashier
	TipsReported          float64        `json:"tips_reported"`
	TotalSellReported     float64        `json:"total_sell_reported"`
	CashIncomesReported   float64        `json:"cash_incomes_reported"`
	OnlineIncomesReported float64        `json:"online_incomes_reported"`
	CardIncomesReported   float64        `json:"card_incomes_reported"`
	CreatedAt             *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt             *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
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
	for _, inv := range invoices {
		total += inv.SubTotal
	}

	return total
}

func GetBruteSell(invoices []invoice.Invoice) float64 {
	total := 0.0
	for _, inv := range invoices {
		total += inv.Total
	}

	return total
}

func GetTotalIncomes(payments []payment.Payment) []TotalIncome {
	totalIncomes := make(map[string]*TotalIncome)
	for _, paymnt := range payments {
		switch paymnt.Method {
		case payment.PaymentMethodCash:
			if _, ok := totalIncomes[payment.PaymentMethodCash]; !ok {
				totalIncomes[payment.PaymentMethodCash] = &TotalIncome{Origin: payment.PaymentMethodCash, Type: "cash", Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCash].Income += paymnt.TotalValue
		case payment.PaymentMethodYuno:
			if _, ok := totalIncomes[payment.PaymentMethodYuno]; !ok {
				totalIncomes[payment.PaymentMethodYuno] = &TotalIncome{Origin: payment.PaymentMethodYuno, Type: "online", Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodYuno].Income += paymnt.TotalValue
		case payment.PaymentMethodCardAmex:
			if _, ok := totalIncomes[payment.PaymentMethodCardAmex]; !ok {
				totalIncomes[payment.PaymentMethodCardAmex] = &TotalIncome{Origin: payment.PaymentMethodCardAmex, Type: "card", Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCardAmex].Income += paymnt.TotalValue
		case payment.PaymentMethodCardMaster:
			if _, ok := totalIncomes[payment.PaymentMethodCardMaster]; !ok {
				totalIncomes[payment.PaymentMethodCardMaster] = &TotalIncome{Origin: payment.PaymentMethodCardMaster, Type: "card", Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCardMaster].Income += paymnt.TotalValue
		case payment.PaymentMethodCardVisa:
			if _, ok := totalIncomes[payment.PaymentMethodCardVisa]; !ok {
				totalIncomes[payment.PaymentMethodCardVisa] = &TotalIncome{Origin: payment.PaymentMethodCardVisa, Type: "card", Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCardVisa].Income += paymnt.TotalValue
		case payment.PaymentMethodCardDinners:
			if _, ok := totalIncomes[payment.PaymentMethodCardDinners]; !ok {
				totalIncomes[payment.PaymentMethodCardDinners] = &TotalIncome{Origin: payment.PaymentMethodCardDinners, Type: "card", Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCardDinners].Income += paymnt.TotalValue
		case payment.PaymentMethodBold:
			if _, ok := totalIncomes[payment.PaymentMethodBold]; !ok {
				totalIncomes[payment.PaymentMethodBold] = &TotalIncome{Origin: payment.PaymentMethodBold, Type: "other", Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodBold].Income += paymnt.TotalValue
		case payment.PaymentMethodBono:
			if _, ok := totalIncomes[payment.PaymentMethodBono]; !ok {
				totalIncomes[payment.PaymentMethodBono] = &TotalIncome{Origin: payment.PaymentMethodBono, Type: "other", Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodBono].Income += paymnt.TotalValue
		}
	}

	totalIncomesSlice := make([]TotalIncome, 0)
	for _, totalIncome := range totalIncomes {
		totalIncomesSlice = append(totalIncomesSlice, *totalIncome)
	}

	return totalIncomesSlice
}

func GetTotalDiscounts(invoices []invoice.Invoice) float64 {
	total := 0.0
	for _, invoice := range invoices {
		total += invoice.TotalDiscounts
	}

	return total
}

func GetTotalSurcharges(invoices []invoice.Invoice) float64 {
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
