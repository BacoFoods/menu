package order

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/BacoFoods/menu/internal"
	accounts "github.com/BacoFoods/menu/pkg/account"
	"github.com/BacoFoods/menu/pkg/client"
	invoices "github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/payment"
	products "github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/shared"
	shifts "github.com/BacoFoods/menu/pkg/shift"
	"github.com/BacoFoods/menu/pkg/tables"
	"github.com/sirupsen/logrus"
)

const (
	LogService string = "pkg/order/service"
)

type Service interface {
	Create(order *Order, ctx context.Context) (*Order, error)
	Update(order *Order) (*Order, error)
	UpdateTable(orderID, tableID uint64) (*Order, error)
	Get(string) (*Order, error)
	Find(filter map[string]any) ([]Order, error)
	UpdateSeats(orderID string, seats int) (*Order, error)
	AddProducts(orderID string, orderItem []OrderItem) (*Order, error)
	RemoveProduct(orderID, productID string) (*Order, error)
	UpdateProduct(product *OrderItem) (*Order, error)
	UpdateStatusNext(orderID string) (*Order, error)
	UpdateStatusPrev(orderID string) (*Order, error)
	ReleaseTable(orderID string) (*Order, error)
	UpdateComments(orderID, comments string) (*Order, error)
	UpdateClientName(orderID, clientName string) (*Order, error)
	AddModifiers(itemID uint, modifiers []OrderModifier) (*OrderItem, error)
	RemoveModifiers(itemID uint, modifiers []OrderModifier) (*OrderItem, error)
	OrderItemUpdateCourse(orderItem *OrderItem) (*OrderItem, error)
	CreateOrderType(orderType *OrderType) (*OrderType, error)
	FindOrderType(filter map[string]any) ([]OrderType, error)
	GetOrderType(orderTypeID string) (*OrderType, error)
	UpdateOrderType(orderTypeID string, orderType *OrderType) (*OrderType, error)
	DeleteOrderType(orderTypeID string) error
	CreateInvoice(orderID string) (*invoices.Invoice, error)
	CalculateInvoice(orderID string, tip *tipData) (*invoices.Invoice, error)
	CalculateInvoiceOIT(orderID string) (*invoices.Invoice, *invoices.Invoice, error)
	Checkout(orderID string, data CheckoutRequest) (*InvoiceCheckout, error)
}

type service struct {
	repository Repository
	table      tables.Repository
	product    products.Repository
	invoice    invoices.Repository
	account    accounts.Repository
	shift      shifts.Repository
	rt         *internal.Rabbit
	payments   payment.Service
}

func NewService(repository Repository,
	table tables.Repository,
	product products.Repository,
	invoice invoices.Repository,
	account accounts.Repository,
	shift shifts.Repository,
	rt *internal.Rabbit,
	payments payment.Service,
) service {
	return service{repository,
		table,
		product,
		invoice,
		account,
		shift,
		rt,
		payments,
	}
}

// Orders
// TODO: improve order creation
func (s service) Create(order *Order, ctx context.Context) (*Order, error) {
	// Setting product items
	productIDs := order.GetProductIDs()
	prods, err := s.product.GetByIDs(productIDs)
	if err != nil {
		shared.LogError("error getting products", LogService, "Create", err, productIDs)
		return nil, fmt.Errorf(ErrorOrderCreation)
	}

	modifierIDs := order.GetModifierIDs()
	modifiers, err := s.product.GetByIDs(modifierIDs)
	if err != nil {
		shared.LogError("error getting modifiers", LogService, "Create", err, modifierIDs)
		return nil, fmt.Errorf(ErrorOrderCreation)
	}

	if len(prods) == 0 && len(productIDs) > 0 {
		return nil, fmt.Errorf(ErrorOrderProductsNotFound)
	}

	order.SetItems(prods, modifiers)
	order.ToInvoice(nil)

	// Setting order status
	order.CurrentStatus = OrderStatusCreated

	// Setting order attendees
	username := ""
	if ctx.Value("account_name") != nil {
		username = ctx.Value("account_name").(string)
	}
	role := ""
	if ctx.Value("role") != nil {
		role = ctx.Value("role").(string)
	}
	accountUUID := ""
	if ctx.Value("account_uuid") != nil {
		accountUUID = ctx.Value("account_uuid").(string)
	}
	channelID := uint(0)
	if value := ctx.Value("channel_id"); value != nil {
		channelIDInt, _ := strconv.Atoi(value.(string))
		channelID = uint(channelIDInt)
	}
	brandID := uint(0)
	if value := ctx.Value("brand_id"); value != nil {
		brandIDInt, _ := strconv.Atoi(value.(string))
		brandID = uint(brandIDInt)
	}
	storeID := uint(0)
	if value := ctx.Value("store_id"); value != nil {
		storeIDInt, _ := strconv.Atoi(value.(string))
		storeID = uint(storeIDInt)
	}
	accountID := uint(0)

	// Setting order shift
	shift, err := s.shift.GetOpenShift(&storeID)
	if err != nil {
		shared.LogWarn("error getting shift", LogService, "Create", err, storeID)
	}
	if shift != nil {
		order.ShiftID = &shift.ID
	}

	newOrder, err := s.repository.Create(order)
	if err != nil {
		shared.LogError("error creating order", LogService, "Create", err, *order)
		return nil, fmt.Errorf(ErrorOrderCreation)
	}

	account, err := s.account.GetByUUID(accountUUID)
	if err != nil {
		shared.LogWarn("error getting account", LogService, "Create", err, accountUUID)
	}

	if account == nil {
		account = &accounts.Account{
			Username:  username,
			ChannelID: &channelID,
			BrandID:   &brandID,
			StoreID:   &storeID,
			Role:      role,
		}
	} else {
		accountID = account.Id
		username = account.DisplayName
		role = account.Role
	}

	attendee := &Attendee{
		AccountID: accountID,
		OrderID:   newOrder.ID,
		Name:      username,
		Role:      role,
		Action:    OrderActionCreated,
		OrderStep: OrderStepCreated,
	}

	if _, err := s.repository.CreateAttendee(attendee); err != nil {
		shared.LogError("error creating attendee", LogService, "Create", err, *attendee)
	}

	// Post the comanda to firebase
	go func() {
		shared.LogInfo(fmt.Sprint(newOrder.Items), LogService, "Create", nil)
		err := s.queueComanda(newOrder.ID, newOrder.TableID, newOrder.StoreID, newOrder.Items)
		if err != nil {
			shared.LogError("error pushing order to firebase", LogService, "Create", err)
		}
	}()

	// Setting table
	// TODO: Send create order and set table to repository to make a trx and rollback if error to avoid has order without table
	if _, err := s.table.SetOrder(newOrder.TableID, &newOrder.ID); err != nil {
		return nil, err
	}

	// Getting order updated from db
	orderDB, err := s.repository.Get(fmt.Sprintf("%d", newOrder.ID))
	if err != nil {
		shared.LogError("error getting order", LogService, "Create", err, newOrder.ID)
		return nil, fmt.Errorf(ErrorOrderCreation)
	}

	return orderDB, nil
}

func (s service) Update(order *Order) (*Order, error) {
	return s.repository.Update(order)
}

func (s service) UpdateTable(orderID, tableID uint64) (*Order, error) {
	order, err := s.repository.Get(fmt.Sprintf("%d", orderID))
	if err != nil {
		shared.LogError("error getting order", LogService, "UpdateTable", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	oldTableID := *order.TableID
	newTableID := uint(tableID)

	if oldTableID == newTableID {
		return order, nil
	}

	if _, err := s.table.SetOrder(&newTableID, &order.ID); err != nil {
		return nil, err
	}

	if _, err := s.table.RemoveOrder(&oldTableID); err != nil {
		return nil, err
	}

	order.TableID = &newTableID
	orderDB, err := s.repository.Update(order)
	if err != nil {
		shared.LogError("error updating order", LogService, "UpdateTable", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	return orderDB, nil
}

func (s service) Get(id string) (*Order, error) {
	return s.repository.Get(id)
}

func (s service) Find(filter map[string]any) ([]Order, error) {
	orders, err := s.repository.Find(filter)
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderFind)
	}

	return orders, nil
}

func (s service) UpdateSeats(orderID string, seats int) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "UpdateSeats", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	if order.Seats == seats {
		return order, nil
	}

	order.Seats = seats
	orderDB, err := s.repository.Update(order)
	if err != nil {
		shared.LogError("error updating order", LogService, "UpdateSeats", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	return orderDB, nil
}

func (s service) AddProducts(orderID string, orderItems []OrderItem) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "AddProduct", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	if order.CurrentStatus != OrderStatusCreated {
		err := fmt.Errorf(ErrorOrderAddProductsForbiddenByStatus)
		shared.LogError("error adding products", LogService, "AddProduct", err, orderID)
		return nil, err
	}

	productIDs := make([]string, len(orderItems))
	modifierIDs := make([]string, 0)
	for i, item := range orderItems {
		if item.ProductID != nil {
			productIDs[i] = fmt.Sprintf("%d", *item.ProductID)
			for _, mod := range item.Modifiers {
				modifierIDs = append(modifierIDs, fmt.Sprintf("%d", *mod.ProductID))
			}
		}
	}

	productsMap, err := s.product.GetAsMapByIDs(productIDs)
	if err != nil {
		shared.LogError("error getting products", LogService, "AddProduct", err, productIDs)
		return nil, fmt.Errorf(ErrorOrderProductGetting)
	}

	productModifiersMap, err := s.product.GetAsMapByIDs(modifierIDs)
	if err != nil {
		shared.LogError("error getting modifiers products", LogService, "AddProduct", err, modifierIDs)
		return nil, fmt.Errorf(ErrorOrderProductGetting)
	}

	newOrderItems := make([]OrderItem, 0)
	errors := ""
	for _, item := range orderItems {
		productID := fmt.Sprintf("%d", *item.ProductID)
		if _, ok := productsMap[productID]; !ok {
			errors += fmt.Sprintf(ErrorOrderProductNotFound, productID)
			continue
		}

		product := productsMap[productID]

		modifiers := make([]OrderModifier, len(item.Modifiers))
		for i, mod := range item.Modifiers {
			productID := fmt.Sprintf("%d", *mod.ProductID)
			if _, ok := productModifiersMap[productID]; !ok {
				errors += fmt.Sprintf(ErrorOrderModifierNotFound, productID)
				continue
			}

			modifier := productModifiersMap[productID]
			modifiers[i] = OrderModifier{
				OrderID:     order.ID,
				ProductID:   mod.ProductID,
				Name:        modifier.Name,
				Description: modifier.Description,
				Image:       modifier.Image,
				SKU:         modifier.SKU,
				Price:       modifier.Price,
				Unit:        modifier.Unit,
				Comments:    mod.Comments,
			}
		}
		newItem := OrderItem{
			OrderID:     &order.ID,
			ProductID:   item.ProductID,
			Name:        product.Name,
			Description: product.Description,
			Image:       product.Image,
			SKU:         product.SKU,
			Price:       product.Price,
			Unit:        product.Unit,
			Comments:    item.Comments,
			Course:      item.Course,
			Modifiers:   modifiers,
		}

		newOrderItems = append(newOrderItems, newItem)
	}

	if errors != "" {
		return nil, fmt.Errorf(errors)
	}

	//	this sets the OrderItem.ID and appends the list to the orignal list of items in the order
	orderDB, err := s.repository.AddProducts(order, newOrderItems)
	if err != nil {
		shared.LogError("error updating order", LogService, "AddProduct", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	if orderDB != nil && len(orderDB.Invoices) != 0 {
		// TODO: improve this to handle multiple invoices
		orderDB.ToInvoice(nil)
	}

	// Post the comanda to firebase
	go func() {
		err := s.queueComanda(order.ID, order.TableID, order.StoreID, newOrderItems)
		if err != nil {
			shared.LogError("error queuing comanda", LogService, "AddProduct", err)
		}
	}()

	return orderDB, nil
}

func (s service) RemoveProduct(orderID, productID string) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "RemoveProduct", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	product, err := s.product.Get(productID)
	if err != nil {
		shared.LogError("error getting product", LogService, "RemoveProduct", err, productID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	order.RemoveProduct(product)

	orderDB, err := s.repository.Update(order)
	if err != nil {
		shared.LogError("error updating order", LogService, "RemoveProduct", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	return orderDB, nil
}

func (s service) UpdateProduct(product *OrderItem) (*Order, error) {
	orderItem, err := s.repository.UpdateOrderItem(product)
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderItemUpdate)
	}

	order, err := s.repository.Get(fmt.Sprintf("%d", orderItem.OrderID))
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	return order, nil
}

func (s service) UpdateStatusNext(orderID string) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "UpdateStatusNext", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	order.UpdateNextStatus()

	if _, err := s.repository.Update(order); err != nil {
		shared.LogError("error updating order status", LogService, "UpdateStatusNext", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdateStatus)
	}

	return order, nil
}

func (s service) UpdateStatusPrev(orderID string) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "UpdateStatusPrev", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	order.UpdatePrevStatus()

	if _, err := s.repository.Update(order); err != nil {
		shared.LogError("error updating order status", LogService, "UpdateStatusPrev", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdateStatus)
	}

	return order, nil
}

func (s service) ReleaseTable(orderID string) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "ReleaseTable", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	tableID := order.TableID
	if order.Table != nil {
		tableID = &order.Table.ID
	}

	if _, err := s.table.RemoveOrder(tableID); err != nil {
		return nil, err
	}

	order.TableID = nil
	order.Table = nil
	orderDB, err := s.repository.Update(order)
	if err != nil {
		shared.LogError("error updating order", LogService, "ReleaseTable", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	return orderDB, nil
}

func (s service) AddModifiers(itemID uint, modifiers []OrderModifier) (*OrderItem, error) {
	orderItem, err := s.repository.GetOrderItem(fmt.Sprintf("%d", itemID))
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderItemGetting)
	}

	orderItem.AddModifiers(modifiers)
	orderItem.SetHash()

	orderItemUpdated, err := s.repository.UpdateOrderItem(orderItem)
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderItemUpdate)
	}

	return orderItemUpdated, nil
}

func (s service) RemoveModifiers(itemID uint, modifiers []OrderModifier) (*OrderItem, error) {
	orderItem, err := s.repository.GetOrderItem(fmt.Sprintf("%d", itemID))
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderItemGetting)
	}

	orderItem.RemoveModifiers(modifiers)
	orderItem.SetHash()

	orderItemUpdated, err := s.repository.UpdateOrderItem(orderItem)
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderItemUpdate)
	}

	return orderItemUpdated, nil
}

func (s service) OrderItemUpdateCourse(orderItem *OrderItem) (*OrderItem, error) {
	orderItem.SetHash()
	orderItemUpdated, err := s.repository.UpdateOrderItem(orderItem)
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderItemUpdate)
	}

	return orderItemUpdated, nil
}

func (s service) UpdateComments(orderID, comments string) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "UpdateComments", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	order.Comments = comments
	return s.repository.Update(order)
}

func (s service) UpdateClientName(orderID, clientName string) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "UpdateClientName", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	order.ClientName = clientName

	return s.repository.Update(order)
}

// Order Types

func (s service) CreateOrderType(orderType *OrderType) (*OrderType, error) {
	return s.repository.CreateOrderType(orderType)
}

func (s service) FindOrderType(filter map[string]any) ([]OrderType, error) {
	return s.repository.FindOrderType(filter)
}

func (s service) GetOrderType(orderTypeID string) (*OrderType, error) {
	return s.repository.GetOrderType(orderTypeID)
}

func (s service) UpdateOrderType(orderTypeID string, orderType *OrderType) (*OrderType, error) {
	return s.repository.UpdateOrderType(orderTypeID, orderType)
}

func (s service) DeleteOrderType(orderTypeID string) error {
	return s.repository.DeleteOrderType(orderTypeID)
}

// Invoice

func (s service) CreateInvoice(orderID string) (*invoices.Invoice, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "CreateInvoice", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	order.ToInvoice(nil)
	if err := order.UpdateStatus(OrderStatusPaying); err != nil {
		shared.LogError("error updating order status", LogService, "CreateInvoice", err, order)
		return nil, fmt.Errorf(ErrorOrderUpdateStatus)
	}
	if _, err := s.repository.Update(order); err != nil {
		shared.LogError("error updating order", LogService, "CreateInvoice", err, order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	invoice := order.Invoices[0]
	invoiceDB, err := s.invoice.CreateUpdate(&invoice)
	if err != nil {
		shared.LogError("error creating invoice", LogService, "CreateInvoice", err, invoice)
		return nil, fmt.Errorf(invoices.ErrorInvoiceCreation)
	}

	return invoiceDB, nil
}

func (s service) CalculateInvoice(orderID string, tip *tipData) (*invoices.Invoice, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "CreateInvoice", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	order.ToInvoice(tip)
	invoice := order.Invoices[0]
	invoice.CalculateTaxDetails()

	return &invoice, nil
}

func (s service) CalculateInvoiceOIT(orderID string) (*invoices.Invoice, *invoices.Invoice, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "CreateInvoice", err, orderID)
		return nil, nil, fmt.Errorf(ErrorOrderGetting)
	}

	dbInvoices, err := s.invoice.Find(map[string]any{"order_id": orderID})
	if err != nil {
		return nil, nil, err
	}

	var oldInvoice *invoices.Invoice
	if len(dbInvoices) > 0 {
		oldInvoice = &dbInvoices[0]
		payments := []payment.Payment{}
		for _, payment := range oldInvoice.Payments {
			if payment.Status != "canceled" {
				payments = append(payments, payment)
			}
		}
		oldInvoice.Payments = payments
	}

	order.ToInvoice(nil)
	newInvoice := order.Invoices[0]
	newInvoice.CalculateTaxDetails()

	return &newInvoice, oldInvoice, nil
}

func (s service) Checkout(orderID string, data CheckoutRequest) (*InvoiceCheckout, error) {
	invoice, err := s.CalculateInvoice(orderID, &tipData{Amount: fmt.Sprint(data.Tip)})
	if err != nil {
		return nil, err
	}

	// invoice.Tip = "percentage"
	// invoice.TipAmount = data.Tip

	// TODO: load client data from ecom
	if data.CustomerID != nil && *data.CustomerID != "" {
		invoice.Client = &client.Client{
			CustomerID: data.CustomerID,
		}
	}

	invoices, err := s.invoice.Find(map[string]any{"order_id": orderID})
	if err != nil {
		return nil, err
	}

	// TODO: asumiendo que solo hay un invoice, con split the bill cambia
	if len(invoices) > 0 {
		oldInvoice := invoices[0]
		invoice.ID = oldInvoice.ID
	}

	// TODO: Se estan creando multiples invoices
	invDB, err := s.invoice.CreateUpdate(invoice)
	if err != nil {
		return nil, err
	}

	// TODO: @Anderson aca se debe pasar el estado de la orden a pagando

	// Paylot immutable
	payment, err := s.payments.CreatePaymentWithPaylot(invDB.ID, invDB.Total, invDB.TipAmount, data.CustomerID)
	if err != nil {
		return nil, err
	}

	return &InvoiceCheckout{
		Payment: payment,
		Invoice: invDB,
	}, nil
}

func (s *service) queueComanda(orderId uint, tableId *uint, storeId *uint, items []OrderItem) error {
	if len(items) == 0 {
		logrus.Info("comanda for order ", orderId, " is empty")
		return nil
	}

	// Timestamp in millis
	logrus.Info("comanda for order ", orderId, " sent")
	ts := time.Now().Unix() * 1000
	data := struct {
		OrderId   uint        `json:"order_id"`
		TableId   *uint       `json:"table_id"`
		Items     []OrderItem `json:"items"`
		Timestamp int64       `json:"timestamp"`
	}{orderId, tableId, items, ts}

	err := s.rt.PutContent(data)

	if err != nil {
		shared.LogError("error queuing comanda", LogService, "queueComanda", err)
	}

	return err
}

var _ Service = &service{}
