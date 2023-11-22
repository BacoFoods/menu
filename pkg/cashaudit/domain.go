package cashaudit

import (
	"github.com/BacoFoods/menu/pkg/invoice"
	orderPKG "github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/payment"
	"github.com/BacoFoods/menu/pkg/shared"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorCashAuditCreating           = "error creating cash audit"
	ErrorCashAuditGetting            = "error getting cash audit"
	ErrorCashAuditGettingStore       = "error getting store"
	ErrorCashAuditGettingStoreID     = "error getting store id"
	ErrorCashAuditGettingLastShift   = "error getting last shift"
	ErrorCashAuditGettingOrders      = "error getting orders"
	ErrorCashAuditNotOrdersFound     = "error not orders found"
	ErrorCashAuditInvalidOrderStatus = "error invalid order status"

	IncomeTypeTip    = "tip"
	IncomeTypeCash   = "cash"
	IncomeTypeOnline = "online"
	IncomeTypeCard   = "card"
	IncomeTypeOther  = "other"

	IncomeDiscrepancyCash = "difference between cash incomes calculated and cash incomes reported"
	IncomeDiscrepancyCard = "difference between card incomes calculated and card incomes reported"
)

type Repository interface {
	Get(storeID string) (*CashAudit, error)
	GetTodayCashAudit(storeID string) (*CashAudit, error)
	Create(cashAudit *CashAudit) (*CashAudit, error)
}

type Income struct {
	ID          uint           `json:"id"`
	CashAuditID *uint          `json:"cash_audit_id"`
	Income      float64        `json:"income" gorm:"precision:18;scale:4"`
	Origin      string         `json:"origin"`
	Type        string         `json:"type" enums:"tip,cash,online,card,other"`
	CreatedAt   *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type CashAudit struct {
	ID                uint       `json:"id"`
	StoreID           *uint      `json:"store_id"`
	StoreName         string     `json:"store_name"`
	ShiftOpen         *time.Time `json:"shift_open"`
	ShiftStartBalance float64    `json:"shift_start_balance"`
	ShiftClose        *time.Time `json:"shift_close"`
	ShiftEndBalance   float64    `json:"shift_end_balance"`
	Orders            uint       `json:"orders"`
	OrdersClosed      uint       `json:"orders_closed"`
	Eaters            uint       `json:"eaters"`
	TotalDiscounts    float64    `json:"discounts"`
	TotalSurcharges   float64    `json:"surcharges"`
	TotalTips         float64    `json:"total_tips"`
	TotalSell         float64    `json:"total_sell"`
	BruteSell         float64    `json:"brute_sell"`
	Incomes           []Income   `json:"total_incomes" gorm:"foreignKey:CashAuditID"`
	// Reported section is for the reported values from cashier
	TipsReported          float64        `json:"tips_reported"`
	TotalSellReported     float64        `json:"total_sell_reported"`
	CashIncomesReported   float64        `json:"cash_incomes_reported"`
	OnlineIncomesReported float64        `json:"online_incomes_reported"`
	CardIncomesReported   float64        `json:"card_incomes_reported"`
	Differences           string         `json:"differences"`  // To save the differences between calculated and reported founded by system
	Observations          string         `json:"observations"` // To save the observations or issues reported from cashier
	CreatedAt             *time.Time     `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt             *time.Time     `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func GetInvoices(orders []orderPKG.Order) []invoice.Invoice {
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

func GetTotalEaters(orders []orderPKG.Order) uint {
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

func GetIncomes(payments []payment.Payment) []Income {
	totalIncomes := make(map[string]*Income)
	for _, paymnt := range payments {
		switch paymnt.Method {
		case payment.PaymentMethodCash:
			if _, ok := totalIncomes[payment.PaymentMethodCash]; !ok {
				totalIncomes[payment.PaymentMethodCash] = &Income{Origin: payment.PaymentMethodCash, Type: IncomeTypeCash, Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCash].Income += paymnt.TotalValue
		case payment.PaymentMethodYuno:
			if _, ok := totalIncomes[payment.PaymentMethodYuno]; !ok {
				totalIncomes[payment.PaymentMethodYuno] = &Income{Origin: payment.PaymentMethodYuno, Type: IncomeTypeOnline, Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodYuno].Income += paymnt.TotalValue
		case payment.PaymentMethodCardAmex:
			if _, ok := totalIncomes[payment.PaymentMethodCardAmex]; !ok {
				totalIncomes[payment.PaymentMethodCardAmex] = &Income{Origin: payment.PaymentMethodCardAmex, Type: IncomeTypeCard, Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCardAmex].Income += paymnt.TotalValue
		case payment.PaymentMethodCardMaster:
			if _, ok := totalIncomes[payment.PaymentMethodCardMaster]; !ok {
				totalIncomes[payment.PaymentMethodCardMaster] = &Income{Origin: payment.PaymentMethodCardMaster, Type: IncomeTypeCard, Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCardMaster].Income += paymnt.TotalValue
		case payment.PaymentMethodCardVisa:
			if _, ok := totalIncomes[payment.PaymentMethodCardVisa]; !ok {
				totalIncomes[payment.PaymentMethodCardVisa] = &Income{Origin: payment.PaymentMethodCardVisa, Type: IncomeTypeCard, Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCardVisa].Income += paymnt.TotalValue
		case payment.PaymentMethodCardDinners:
			if _, ok := totalIncomes[payment.PaymentMethodCardDinners]; !ok {
				totalIncomes[payment.PaymentMethodCardDinners] = &Income{Origin: payment.PaymentMethodCardDinners, Type: IncomeTypeCard, Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodCardDinners].Income += paymnt.TotalValue
		case payment.PaymentMethodBold:
			if _, ok := totalIncomes[payment.PaymentMethodBold]; !ok {
				totalIncomes[payment.PaymentMethodBold] = &Income{Origin: payment.PaymentMethodBold, Type: IncomeTypeOther, Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodBold].Income += paymnt.TotalValue
		case payment.PaymentMethodBono:
			if _, ok := totalIncomes[payment.PaymentMethodBono]; !ok {
				totalIncomes[payment.PaymentMethodBono] = &Income{Origin: payment.PaymentMethodBono, Type: IncomeTypeOther, Income: paymnt.TotalValue}
				continue
			}
			totalIncomes[payment.PaymentMethodBono].Income += paymnt.TotalValue
		default:
			shared.LogWarn("payment method not tracked", LogService, "GetIncomes", nil, paymnt)
			if _, ok := totalIncomes[payment.PaymentMethodUntracked]; !ok {
				totalIncomes[payment.PaymentMethodUntracked] = &Income{Origin: payment.PaymentMethodUntracked, Type: IncomeTypeOther, Income: paymnt.TotalValue}
				continue
			}
		}
	}

	totalIncomesSlice := make([]Income, 0)
	for _, totalIncome := range totalIncomes {
		totalIncomesSlice = append(totalIncomesSlice, *totalIncome)
	}

	return totalIncomesSlice
}

func GetTipIncomes(payments []payment.Payment) []Income {
	totalTips := make(map[string]*Income)

	for _, paymnt := range payments {
		if paymnt.Tip > 0 {
			if _, ok := totalTips[payment.PaymentMethodCash]; !ok {
				totalTips[payment.PaymentMethodCash] = &Income{Origin: payment.PaymentMethodCash, Type: IncomeTypeTip, Income: paymnt.Tip}
				continue
			}
			totalTips[payment.PaymentMethodCash].Income += paymnt.Tip
		}
	}

	totalTipsSlice := make([]Income, 0)
	for _, totalTip := range totalTips {
		totalTipsSlice = append(totalTipsSlice, *totalTip)
	}
	return totalTipsSlice
}

func GetTotalTips(payments []payment.Payment) float64 {
	total := 0.0
	for _, payment := range payments {
		total += payment.Tip
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

func GetOrdersClosed(orders []orderPKG.Order) uint {
	closed := 0
	for _, order := range orders {
		if order.CurrentStatus == orderPKG.OrderStatusClosed {
			closed++
		}
	}
	return uint(closed)
}
