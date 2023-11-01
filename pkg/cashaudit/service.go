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
	lastShift, err := s.shifts.GetLastShift(storeID)
	if err != nil {
		shared.LogError("error getting last shift", LogService, "Get", err, storeID)
		return nil, fmt.Errorf(ErrorCashAuditGettingLastShift)
	}
	cashAudit.ShiftOpen = lastShift.StartTime
	cashAudit.ShiftStartBalance = lastShift.StartBalance
	cashAudit.ShiftClose = lastShift.EndTime
	cashAudit.ShiftEndBalance = lastShift.EndBalance

	// Set Orders details
	orderList, err := s.orders.FindByShift(lastShift.ID)
	if err != nil {
		shared.LogError("error getting orders by shift", LogService, "Get", err, lastShift.ID)
		return nil, fmt.Errorf(ErrorCashAuditGettingOrders)
	}
	cashAudit.Orders = uint(len(orderList))
	cashAudit.Eaters = GetTotalEaters(orderList)

	// Set Invoices details
	invoiceList := GetInvoices(orderList)
	cashAudit.Discounts = GetTotalDiscounts(invoiceList)
	cashAudit.Surcharges = GetTotalSourcharges(invoiceList)
	cashAudit.Tips = GetTotalTips(invoiceList)
	cashAudit.TotalSell = GetTotalSell(invoiceList)

	return &cashAudit, nil
}