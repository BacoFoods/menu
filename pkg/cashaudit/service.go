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
}

type service struct {
	stores   store.Repository
	orders   order.Repository
	invoices invoice.Repository
	shifts   shift.Repository
}

func NewService(stores store.Repository,
	orders order.Repository,
	invoices invoice.Repository,
	shifts shift.Repository) service {
	return service{
		stores,
		orders,
		invoices,
		shifts,
	}
}

func (s service) Get(storeID string) (*CashAudit, error) {
	cashAudit := CashAudit{}

	// Set Store details
	auditStore, err := s.stores.Get(storeID)
	if err != nil {
		shared.LogError("error getting store", LogService, "Get", err, storeID)
		return nil, fmt.Errorf(ErrorCashAuditGettingStore)
	}
	cashAudit.StoreName = auditStore.Name

	// Set Shift details
	// TODO: Shifts not already implemented, so create shifts to work with balance

	// Set Orders details
	orderList, err := s.orders.GetLastDayOrders(storeID)
	if err != nil {
		shared.LogError("error getting last day orders", LogService, "Get", err, storeID)
		return nil, fmt.Errorf(ErrorCashAuditGettingOrders)
	}

	if len(orderList) == 0 {
		return nil, fmt.Errorf(ErrorCashAuditNotOrdersFound)
	}

	cashAudit.Orders = uint(len(orderList))
	cashAudit.Eaters = GetTotalEaters(orderList)

	// Getting invoices from orders
	invoiceList := GetInvoices(orderList)
	paymentsList := GetPayments(invoiceList)

	cashAudit.Orders = uint(len(orderList))
	cashAudit.Eaters = GetTotalEaters(orderList)
	cashAudit.TotalIncomes = GetTotalIncomes(paymentsList)
	cashAudit.TotalSell = GetTotalSell(invoiceList)
	cashAudit.BruteSell = GetBruteSell(invoiceList)

	return &cashAudit, nil
}

func (s service) Create(storeID string, cashReported *CashAudit) (*CashAudit, error) {
	// Set Store details
	auditStore, err := s.stores.Get(storeID)
	if err != nil {
		shared.LogError("error getting store", LogService, "Get", err, storeID)
		return nil, fmt.Errorf(ErrorCashAuditGettingStore)
	}
	cashReported.StoreName = auditStore.Name

	// Getting day's orders
	orderList, err := s.orders.GetLastDayOrders(storeID)
	if err != nil {
		return nil, fmt.Errorf(ErrorCashAuditGettingOrders)
	}

	if len(orderList) == 0 {
		return nil, fmt.Errorf(ErrorCashAuditNotOrdersFound)
	}

	// Getting invoices from orders
	invoiceList := GetInvoices(orderList)
	paymentsList := GetPayments(invoiceList)

	cashReported.Orders = uint(len(orderList))
	cashReported.Eaters = GetTotalEaters(orderList)
	cashReported.TotalIncomes = GetTotalIncomes(paymentsList)
	cashReported.TotalSell = GetTotalSell(invoiceList)
	cashReported.BruteSell = GetBruteSell(invoiceList)

	return cashReported, nil
}
