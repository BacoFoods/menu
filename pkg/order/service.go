package order

import (
	"context"
	"fmt"
	accounts "github.com/BacoFoods/menu/pkg/account"
	invoices "github.com/BacoFoods/menu/pkg/invoice"
	products "github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/shared"
	statuses "github.com/BacoFoods/menu/pkg/status"
	"github.com/BacoFoods/menu/pkg/tables"
	"strconv"
)

const (
	LogService string = "pkg/order/service"
)

type Service interface {
	Create(order *Order, ctx context.Context) (*Order, error)
	UpdateTable(orderID, tableID uint64) (*Order, error)
	Get(string) (*Order, error)
	Find(filter map[string]any) ([]Order, error)
	UpdateSeats(orderID string, seats int) (*Order, error)
	AddProducts(orderID string, orderItem []OrderItem) (*Order, error)
	RemoveProduct(orderID, productID string) (*Order, error)
	UpdateProduct(product *OrderItem) (*Order, error)
	UpdateStatus(orderID, statusCode string) (*Order, error)
	ReleaseTable(orderID string) (*Order, error)
	UpdateComments(orderID, comments string) (*Order, error)
	AddModifiers(itemID uint, modifiers []OrderModifier) (*OrderItem, error)
	RemoveModifiers(itemID uint, modifiers []OrderModifier) (*OrderItem, error)
	OrderItemUpdateCourse(orderItem *OrderItem) (*OrderItem, error)
	CreateOrderType(orderType *OrderType) (*OrderType, error)
	FindOrderType(filter map[string]any) ([]OrderType, error)
	GetOrderType(orderTypeID string) (*OrderType, error)
	UpdateOrderType(orderTypeID string, orderType *OrderType) (*OrderType, error)
	DeleteOrderType(orderTypeID string) error
	CreateInvoice(orderID string) (*invoices.Invoice, error)
}

type service struct {
	repository Repository
	table      tables.Repository
	product    products.Repository
	invoice    invoices.Repository
	status     statuses.Repository
	account    accounts.Repository
}

func NewService(repository Repository,
	table tables.Repository,
	product products.Repository,
	invoice invoices.Repository,
	status statuses.Repository,
	account accounts.Repository) service {
	return service{repository,
		table,
		product,
		invoice,
		status,
		account,
	}
}

// Orders

func (s service) Create(order *Order, ctx context.Context) (*Order, error) {
	// Setting product items
	productIDs := order.GetProductIDs()

	prods, err := s.product.GetByIDs(productIDs)
	if err != nil {
		shared.LogError("error getting products", LogService, "Create", err, productIDs)
		return nil, fmt.Errorf(ErrorOrderCreation)
	}

	if len(prods) == 0 {
		return nil, fmt.Errorf(ErrorOrderProductsNotFound)
	}

	order.SetItems(prods)

	newOrder, err := s.repository.Create(order)
	if err != nil {
		shared.LogError("error creating order", LogService, "Create", err, *order)
		return nil, fmt.Errorf(ErrorOrderCreation)
	}

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
	channelID := int64(0)
	if value := ctx.Value("channel_id"); value != nil {
		channelID, _ = strconv.ParseInt(value.(string), 10, 64)
	}
	brandID := int64(0)
	if value := ctx.Value("brand_id"); value != nil {
		brandID, _ = strconv.ParseInt(value.(string), 10, 64)
	}
	storeID := int64(0)
	if value := ctx.Value("store_id"); value != nil {
		storeID, _ = strconv.ParseInt(value.(string), 10, 64)
	}
	accountID := uint(0)

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

	// Setting table
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
	return s.repository.Find(filter)
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

		order.AddProduct(OrderItem{
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
		})
	}

	if errors != "" {
		return nil, fmt.Errorf(errors)
	}

	orderDB, err := s.repository.Update(order)
	if err != nil {
		shared.LogError("error updating order", LogService, "AddProduct", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

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

func (s service) UpdateStatus(orderID, statusCode string) (*Order, error) {
	order, err := s.repository.Get(orderID)
	if err != nil {
		shared.LogError("error getting order", LogService, "UpdateStatus", err, orderID)
		return nil, fmt.Errorf(ErrorOrderGetting)
	}

	status, err := s.status.GetByCode(statusCode)
	if err != nil {
		shared.LogError("error getting status", LogService, "UpdateStatus", err, status)
		return nil, fmt.Errorf(ErrorOrderGettingStatus)
	}

	if err := order.UpdateStatus(status); err != nil {
		return nil, err
	}

	if _, err := s.repository.Update(order); err != nil {
		shared.LogError("error updating order", LogService, "UpdateStatus", err, *order)
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

	if _, err := s.table.RemoveOrder(order.TableID); err != nil {
		return nil, err
	}

	order.TableID = nil
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

	orderItemUpdated, err := s.repository.UpdateOrderItem(orderItem)
	if err != nil {
		return nil, fmt.Errorf(ErrorOrderItemUpdate)
	}

	return orderItemUpdated, nil
}

func (s service) OrderItemUpdateCourse(orderItem *OrderItem) (*OrderItem, error) {
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

	order.ToInvoice()
	invoice, err := s.invoice.CreateUpdate(order.Invoice)
	if err != nil {
		shared.LogError("error creating invoice", LogService, "CreateInvoice", err, order.Invoice)
		return nil, fmt.Errorf(invoices.ErrorInvoiceCreation)
	}

	if order.InvoiceID != nil {
		return invoice, nil
	}

	order.InvoiceID = &invoice.ID
	if _, err = s.repository.Update(order); err != nil {
		shared.LogError("error updating order", LogService, "CreateInvoice", err, *order)
		return nil, fmt.Errorf(ErrorOrderUpdate)
	}

	return invoice, nil
}
