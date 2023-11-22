package cashaudit

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/BacoFoods/menu/pkg/shift"
	"github.com/BacoFoods/menu/pkg/store"
)

const (
	LogService = "pkg/cashaudit/service"
)

type Service interface {
	Get(storeID string) (*CashAudit, error)
	Create(storeID string, cashAudit *CashAudit) (*CashAudit, error)
	Confirm(cashAuditID string) (*CashAudit, error)
}

type service struct {
	repository Repository
	stores     store.Repository
	orders     order.Repository
	invoices   invoice.Repository
	shifts     shift.Repository
}

func NewService(repository Repository,
	stores store.Repository,
	orders order.Repository,
	invoices invoice.Repository,
	shifts shift.Repository) service {
	return service{
		repository,
		stores,
		orders,
		invoices,
		shifts,
	}
}

func (s service) Get(storeID string) (*CashAudit, error) {
	return s.calculateCashAudit(storeID, order.OrderStatusClosed)
}

func (s service) Create(storeID string, cashReported *CashAudit) (*CashAudit, error) {
	todayCashAudit, err := s.repository.GetTodayCashAudit(storeID)
	if err != nil {
		return nil, err
	}

	if todayCashAudit != nil {
		return todayCashAudit, nil
	}

	// Set Store details
	cashAudit, err := s.calculateCashAudit(storeID, order.OrderStatusClosed)
	if err != nil {
		return nil, err
	}

	cashAudit.TipsReported = cashReported.TipsReported
	cashAudit.TotalSellReported = cashReported.TotalSellReported
	cashAudit.CashIncomesReported = cashReported.CashIncomesReported
	cashAudit.OnlineIncomesReported = cashReported.OnlineIncomesReported
	cashAudit.CardIncomesReported = cashReported.CardIncomesReported

	// Validate discrepancies between reported and calculated
	cashAudit.Differences = s.validateDiscrepancy(cashAudit)

	// Create cash audit
	cashAudit, err = s.repository.Create(cashAudit)
	if err != nil {
		return nil, err
	}

	return cashAudit, nil
}

func (s service) Confirm(cashAuditID string) (*CashAudit, error) {
	cashAudit, err := s.repository.Get(cashAuditID)
	if err != nil {
		return nil, err
	}

	cashAudit.Confirmation = true

	cashAudit, err = s.repository.Update(cashAudit)
	if err != nil {
		return nil, err
	}

	return cashAudit, nil
}

func (s service) calculateCashAudit(storeID string, status string) (*CashAudit, error) {
	cashAudit := CashAudit{}

	// Validate status
	if !order.OrderStatusValid(status) {
		return nil, fmt.Errorf(ErrorCashAuditInvalidOrderStatus)
	}

	// Set Store details
	auditStore, err := s.stores.Get(storeID)
	if err != nil {
		shared.LogError("error getting store", LogService, "Get", err, storeID)
		return nil, fmt.Errorf(ErrorCashAuditGettingStore)
	}
	cashAudit.StoreID = &auditStore.ID
	cashAudit.StoreName = auditStore.Name

	// Getting day's orders by status
	orderList, err := s.orders.GetLastDayOrdersByStatus(storeID, status)
	if err != nil {
		return nil, fmt.Errorf(ErrorCashAuditGettingOrders)
	}

	if len(orderList) == 0 {
		return nil, fmt.Errorf(ErrorCashAuditNotOrdersFound)
	}

	// Getting invoices from orders
	invoiceList := GetInvoices(orderList)
	paymentsList := GetPayments(invoiceList)

	cashAudit.Orders = uint(len(orderList))
	cashAudit.OrdersClosed = GetOrdersClosed(orderList)
	cashAudit.Eaters = GetTotalEaters(orderList)

	// Setting incomes
	cashAudit.Incomes = GetIncomes(paymentsList)

	// Setting tips
	cashAudit.Incomes = append(cashAudit.Incomes, GetTipIncomes(paymentsList)...)

	cashAudit.TotalTips = GetTotalTips(paymentsList)
	cashAudit.TotalDiscounts = GetTotalDiscounts(invoiceList)
	cashAudit.TotalSell = GetTotalSell(invoiceList)
	cashAudit.BruteSell = GetBruteSell(invoiceList)

	return &cashAudit, nil
}

func (s service) validateDiscrepancy(cashAudit *CashAudit) string {
	cashIncomesTotal := 0.0
	for _, income := range cashAudit.Incomes {
		if income.Type == IncomeTypeCash {
			cashIncomesTotal += income.Income
		}
	}

	if cashAudit.CashIncomesReported != cashIncomesTotal {
		return IncomeDiscrepancyCash
	}

	cardIncomesTotal := 0.0
	for _, income := range cashAudit.Incomes {
		if income.Type == IncomeTypeCard {
			cardIncomesTotal += income.Income
		}
	}

	if cashAudit.CardIncomesReported != cardIncomesTotal {
		return IncomeDiscrepancyCard
	}

	return ""
}
