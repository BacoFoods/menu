package order

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/BacoFoods/menu/pkg/channel"
	discountPKG "github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/payment"
	"github.com/BacoFoods/menu/pkg/shared"

	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/client"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/tables"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	ErrorBadRequest                        = "error bad request"
	ErrorBadRequestOrderID                 = "error bad request wrong order id"
	ErrorBadRequestOrderItemID             = "error bad request wrong order item id"
	ErrorBadRequestProductID               = "error bad request wrong product id"
	ErrorBadRequestTableID                 = "error bad request wrong table id"
	ErrorBadRequestStoreID                 = "error bad request wrong store id"
	ErrorBadRequestOrderSeats              = "error bad request wrong order seats can't be less than 0"
	ErrorOrderDeleting                     = "error deleting order"
	ErrorOrderCreation                     = "error creating order"
	ErrorOrderGetting                      = "error getting order"
	ErrorOrderFind                         = "error finding orders"
	ErrorOrderAddProductsForbiddenByStatus = "error adding products to order forbidden by status"
	ErrorOrderUpdate                       = "error updating order"
	ErrorOrderUpdateStatus                 = "error updating order status"
	ErrorOrderUpdateInvalidStatus          = "error updating order invalid status"
	ErrorOrderProductGetting               = "error getting order product"
	ErrorOrderProductNotFound              = "error order product with id %v not found; "
	ErrorOrderProductsNotFound             = "error order products not found"
	ErrorOrderModifierNotFound             = "error order modifier with id %v not found; "
	ErrorOrderUpdatingComments             = "error updating order comments"
	ErrorOrderUpdatingClientName           = "error updating order client name"
	ErrorOrderUpdatingStatus               = "error updating order status"
	ErrorOrderInvoiceCreation              = "error creating order invoice"
	ErrorOrderInvoiceUpdate                = "error updating order invoice"
	ErrorOrderInvoiceCreationDiscounts     = "error creating order invoice getting discounts"
	ErrorOrderInvoiceGettingClient         = "error getting order invoice client"
	ErrorOrderInvoiceCalculation           = "error calculating invoice"
	ErrorOrderClosed                       = "error order is closed"
	ErrorOrderIDEmpty                      = "order id is empty"

	ErrorOrderItemUpdate       = "error updating order item"
	ErrorOrderItemGetting      = "error getting order item"
	ErrorOrderItemUpdateCourse = "error updating order item course"

	ErrorOrderTypeCreation               = "error creating order type"
	ErrorOrderTypeFinding                = "error finding order type"
	ErrorOrderTypeGetting                = "error getting order type"
	ErrorOrderTypeUpdating               = "error updating order type"
	ErrorOrderTypeDeleting               = "error deleting order type"
	ErrorOrderInvoicePlemsiBuilding      = "error building plemsi invoice"
	ErrorOrderInvoiceFacturacionConfig   = "error getting facturacion config"
	ErrorOrderInvoiceInvalidDocumentType = "error invalid document type"
	ErrorOrderInvoiceEmission            = "error emitting order invoice"

	TaxPercentage = 0.08

	OrderStepCreated  OrderStep = "created"
	OrderStepClosed   OrderStep = "closed"
	OrderStepInvoiced OrderStep = "invoiced"

	OrderActionCreated  OrderAction = "fue atendido por"
	OrderActionClosed   OrderAction = "fue cerrada por"
	OrderActionInvoiced OrderAction = "fue facturado por"

	LogDomain = "pkg/order/domain"

	OrderStatusCreated = "created"
	OrderStatusPaying  = "paying"
	OrderStatusClosed  = "closed"
)

func applyDiscount(value float64, discounts []invoice.DiscountApplied) (newValue float64, appliedDiscount float64) {
	newValue = value

	for _, discount := range discounts {
		newValue = discount.ApplyRounded(newValue)
	}

	newValue = math.Round(newValue)
	appliedDiscount = value - newValue

	return newValue, appliedDiscount
}

func OrderStatusValid(status string) bool {
	switch status {
	case OrderStatusCreated, OrderStatusPaying, OrderStatusClosed:
		return true
	default:
		return false
	}
}

type OrderStep string

type OrderAction string

type Repository interface {
	// Order
	Create(order *Order, ch *channel.Channel) (*Order, error)
	Get(orderID string) (*Order, error)
	Find(filter map[string]any) ([]Order, error)
	Update(order *Order) (*Order, error)
	FindByShift(shiftID uint) ([]Order, error)
	AddProducts(order *Order, newItems []OrderItem) (*Order, error)
	GetLastDayOrders(storeID string) ([]Order, error)
	GetLastDayOrdersByStatus(storeID string, status string) ([]Order, error)
	Delete(orderID string) error
	FindByIdempotencyKey(idempotencyKey string, storeID *uint) (*Order, error)

	// OrderItem
	UpdateOrderItem(orderItem *OrderItem) (*OrderItem, error)
	GetOrderItem(orderItemID string) (*OrderItem, error)
	UpdateTable(order *Order, newTableID uint) (*Order, error)

	// OrderType
	CreateOrderType(*OrderType) (*OrderType, error)
	FindOrderType(filter map[string]any) ([]OrderType, error)
	GetOrderType(orderTypeID string) (*OrderType, error)
	UpdateOrderType(orderTypeID string, orderType *OrderType) (*OrderType, error)
	DeleteOrderType(orderTypeID string) error

	// Attendee
	CreateAttendee(attendee *Attendee) (*Attendee, error)
}

type Order struct {
	ID             uint              `json:"id" gorm:"primaryKey"`
	Statuses       []OrderStatus     `json:"status" gorm:"foreignKey:OrderID" swaggerignore:"true"`
	Code           string            `json:"code" swaggerignore:"true"`
	CurrentStatus  string            `json:"current_status" gorm:"index:idx_orders_current_status,default:'created'"`
	OrderType      string            `json:"order_type"`
	ClientName     string            `json:"client_name"`
	BrandID        *uint             `json:"brand_id" binding:"required"`
	Brand          *brand.Brand      `json:"brand,omitempty" swaggerignore:"true"`
	StoreID        *uint             `json:"store_id" binding:"required" gorm:"index:idx_orders_store_id"`
	Store          *store.Store      `json:"store,omitempty" swaggerignore:"true"`
	ChannelID      *uint             `json:"channel_id" binding:"required"`
	TableID        *uint             `json:"table_id"`
	Table          *tables.Table     `json:"table" swaggerignore:"true" gorm:"foreignKey:TableID"`
	TypeID         *uint             `json:"type_id"`
	Type           *OrderType        `json:"type"`
	Comments       string            `json:"comments"`
	Items          []OrderItem       `json:"items"  gorm:"foreignKey:OrderID"`
	CookingTime    int               `json:"cooking_time"`
	Seats          int               `json:"seats"`
	ExternalCode   string            `json:"external_code"`
	Invoices       []invoice.Invoice `json:"invoices"  gorm:"foreignKey:OrderID" swaggerignore:"true"`
	Attendees      []Attendee        `json:"attendees" gorm:"foreignKey:OrderID"`
	ShiftID        *uint             `json:"shift_id"`
	IdempotencyKey *string           `json:"idempotency_key"`
	ClosedAt       *time.Time        `json:"closed_at" swaggerignore:"true"`
	CreatedAt      *time.Time        `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt      *time.Time        `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt      *gorm.DeletedAt   `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (o *Order) GetProductIDs() []string {
	ids := make([]string, len(o.Items))
	for i, item := range o.Items {
		ids[i] = fmt.Sprintf("%d", *item.ProductID)
	}
	return ids
}

func (o *Order) GetModifierIDs() []string {
	ids := make([]string, 0)
	for _, item := range o.Items {
		for _, modifier := range item.Modifiers {
			ids = append(ids, fmt.Sprintf("%d", *modifier.ProductID))
		}
	}
	return ids
}

func (o *Order) SetItems(products []product.Product, modifiers []product.Product) {
	productsMap := make(map[string]product.Product)
	for _, p := range products {
		productsMap[fmt.Sprintf("%d", p.ID)] = p
	}

	modifiersMap := make(map[string]product.Product)
	for _, m := range modifiers {
		modifiersMap[fmt.Sprintf("%d", m.ID)] = m
	}

	items := make([]OrderItem, 0)
	for _, item := range o.Items {
		if p, ok := productsMap[fmt.Sprintf("%d", *item.ProductID)]; ok {
			item.Name = p.Name
			item.Description = p.Description
			item.Image = p.Image
			item.SKU = p.SKU
			item.Price = p.Price
			item.Unit = p.Unit

			if p.TaxBase == 0 && p.Tax != nil {
				item.Tax = p.Tax.Name
				item.TaxPercentage = p.Tax.Percentage
				item.TaxBase = p.Price / (1 + p.Tax.Percentage)
				item.TaxAmount = p.Price - item.TaxBase
			}

			item.SetHash()

			modifierList := make([]OrderModifier, 0)
			for _, modifier := range item.Modifiers {
				if m, ok := modifiersMap[fmt.Sprintf("%d", *modifier.ProductID)]; ok {
					modifier.Name = m.Name
					modifier.Description = m.Description
					modifier.Image = m.Image
					modifier.SKU = m.SKU
					modifier.Price = m.Price
					modifier.Unit = m.Unit
					modifier.ProductID = &m.ID
					modifier.OrderID = o.ID
					modifierList = append(modifierList, modifier)
				}
			}
			item.Modifiers = modifierList
			items = append(items, item)
		}
	}
	o.Items = items
}

func (o *Order) AddProduct(orderItem OrderItem) {
	o.Items = append(o.Items, orderItem)
}

func (o *Order) RemoveProduct(product *product.Product) {
	for i, item := range o.Items {
		if *item.ProductID == product.ID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			return
		}
	}
}

// ToInvoice uses next definitions:
// Product Price: is the price of the product without any discount and taxes included
// Product Discounted Price: is the price of the product after applying discounts, discount is applied to the product price
// Product Base Tax: is the tax base of the product, tax base is the price of the product without taxes
// Product Tax: is the tax amount of the product
// TotalTips: is the sum of all tips
// SubTotal: is the sum of all Product Prices
// BaseTax: is the sum of all Product Base Taxes
// Total: is the sum of SubTotal(taxes are included) + TotalTips - TotalDiscounts
func (o *Order) ToInvoice(tip *TipData, discounts ...discountPKG.Discount) {
	// Remove invoices
	o.Invoices = nil
	subtotal := 0.0
	newInvoice := invoice.Invoice{
		OrderID:   &o.ID,
		BrandID:   o.BrandID,
		StoreID:   o.StoreID,
		ChannelID: o.ChannelID,
		TableID:   o.TableID,
		ShiftID:   o.ShiftID,
		Items:     make([]invoice.Item, 0),
		Discounts: make([]invoice.DiscountApplied, 0),
		Client:    client.DefaultClient(),
		BaseTax:   0,
	}

	// Adding discounts to invoice
	for _, d := range discounts {
		if d.Type != discountPKG.DiscountTypePercentage {
			// TODO: we only support percentage discounts as applying a value discounts
			// makes tax calculations more complex
			continue
		}

		newInvoice.Discounts = append(newInvoice.Discounts, invoice.DiscountApplied{
			DiscountID:  d.ID,
			Name:        d.Name,
			Description: d.Description,
			Percentage:  d.Percentage,
			Amount:      0,
			Type:        string(d.Type),
		})
	}

	// Adding items to invoice
	orderItems := make([]OrderItem, 0)
	for _, orderItem := range o.Items {

		// Default tax values
		if orderItem.Tax == "" {
			orderItem.Tax = "ico" // Default tax type
		}

		if orderItem.TaxPercentage == 0 {
			orderItem.TaxPercentage = 0.08 // Default tax percentage
		}

		if orderItem.TaxBase == 0 {
			orderItem.TaxBase = math.Ceil(orderItem.Price / (1 + orderItem.TaxPercentage)) // Default tax base
		}

		if orderItem.TaxAmount == 0 {
			orderItem.TaxAmount = orderItem.Price - orderItem.TaxBase // Default tax amount
		}

		// Discounts
		orderItem.DiscountedPrice = orderItem.Price
		orderItem.Discount = 0

		for _, discount := range newInvoice.Discounts {
			orderItem.DiscountedPrice -= discount.CalculateAmountRounded(orderItem.Price)
			orderItem.DiscountPercent += discount.Percentage
			orderItem.Discount += discount.CalculateAmountRounded(orderItem.Price)
			orderItem.DiscountReason += fmt.Sprintf("%s - %s - %.2f - %.2f - applied to: %.2f;", discount.Name, discount.Description, discount.Percentage, discount.CalculateAmountRounded(orderItem.Price), orderItem.Price)
		}

		newInvoice.BaseTax += orderItem.TaxBase
		newInvoice.Taxes += orderItem.TaxAmount
		newInvoice.TotalDiscounts += orderItem.Discount

		newInvoice.Items = append(newInvoice.Items, invoice.Item{
			ProductID:          orderItem.ProductID,
			Name:               orderItem.Name,
			Description:        orderItem.Description,
			SKU:                orderItem.SKU,
			Price:              orderItem.Price,
			Comments:           orderItem.Comments,
			Hash:               orderItem.Hash,
			DiscountedPrice:    orderItem.DiscountedPrice,
			DiscountPercentage: orderItem.DiscountPercent,
			DiscountReason:     orderItem.DiscountReason,
			DiscountAmount:     orderItem.Discount,
			Tax:                orderItem.Tax,
			TaxPercentage:      orderItem.TaxPercentage,
			TaxAmount:          orderItem.TaxAmount,
			TaxBase:            orderItem.TaxBase,
		})

		// Adding orderItem price to subtotal
		subtotal += orderItem.Price

		for _, modifier := range orderItem.Modifiers {

			// Discounts
			modifier.DiscountedPrice = orderItem.Price
			orderItem.Discount = 0

			for _, discount := range newInvoice.Discounts {
				modifier.DiscountedPrice -= discount.CalculateAmountRounded(orderItem.Price)
				orderItem.DiscountPercent += discount.Percentage
				orderItem.Discount += discount.CalculateAmountRounded(orderItem.Price)
				orderItem.DiscountReason += fmt.Sprintf("%s - %s - %.2f - %.2f - applied to: %.2f;", discount.Name, discount.Description, discount.Percentage, discount.CalculateAmountRounded(orderItem.Price), orderItem.Price)
			}

			// Adding modifier price to subtotal
			modifierPrice, appliedDiscount := applyDiscount(modifier.Price, newInvoice.Discounts)
			subtotal += modifierPrice

			modifierTax := "ico" // Default tax
			if modifier.Tax != "" {
				modifierTax = modifier.Tax
			}

			modifierTaxPerc := 0.08 // TODO Improve tax calculation
			if modifier.TaxPercentage != 0 {
				modifierTaxPerc = modifier.TaxPercentage
			}

			modifierBaseTax := math.Floor(modifierPrice / (1 + modifierTaxPerc))
			newInvoice.BaseTax += modifierBaseTax
			newInvoice.Taxes += modifierPrice - modifierBaseTax
			newInvoice.TotalDiscounts += appliedDiscount

			newInvoice.Items = append(newInvoice.Items, invoice.Item{
				ProductID:       modifier.ProductID,
				Name:            modifier.Name,
				Description:     modifier.Description,
				SKU:             modifier.SKU,
				Price:           modifier.Price,
				Comments:        modifier.Comments,
				DiscountedPrice: modifierPrice,
				Tax:             modifierTax,
				TaxPercentage:   modifierTaxPerc,
			})
		}

		orderItems = append(orderItems, orderItem)
	}

	// Setting order items updates
	o.Items = orderItems

	// Setting subtotals, subtotals includes taxes
	newInvoice.SubTotal = subtotal

	// Setting taxes
	newInvoice.CalculateTaxDetails()

	// Setting tips
	if tip != nil {
		tipType, tipValue := tip.GetValueAndType()
		newInvoice.Tip = fmt.Sprintf("%s - %f", tipType, tipValue)
		if tipType == "percentage" {
			tipAmount := math.Floor(newInvoice.BaseTax * tipValue)
			newInvoice.TipAmount = tipAmount
		} else {
			newInvoice.TipAmount = tipValue
		}
	}

	newInvoice.Total = newInvoice.SubTotal + newInvoice.TipAmount - newInvoice.TotalDiscounts

	// Setting invoice
	o.Invoices = []invoice.Invoice{newInvoice}
}

func (o *Order) UpdateStatus(status string) error {
	if status == o.CurrentStatus {
		shared.LogWarn("order already has this status", LogDomain, "UpdateStatus", nil, o.ID, o.CurrentStatus, status)
		return nil
	}

	if status == OrderStatusCreated && strings.TrimSpace(o.CurrentStatus) == "" {
		o.CurrentStatus = status
		o.Statuses = append(o.Statuses, OrderStatus{
			Code:    OrderStatusCreated,
			OrderID: &o.ID,
		})
		return nil
	}

	if status == OrderStatusPaying && o.CurrentStatus == OrderStatusCreated {
		o.CurrentStatus = status
		o.Statuses = append(o.Statuses, OrderStatus{
			Code:    OrderStatusPaying,
			OrderID: &o.ID,
		})
		return nil
	}

	if status == OrderStatusClosed && o.CurrentStatus == OrderStatusPaying {
		o.CurrentStatus = status
		o.Statuses = append(o.Statuses, OrderStatus{
			Code:    OrderStatusClosed,
			OrderID: &o.ID,
		})
		return nil
	}

	err := fmt.Errorf(ErrorOrderUpdateStatus)
	shared.LogError(ErrorOrderUpdateStatus, LogDomain, "UpdateStatus", nil, o.ID, o.CurrentStatus, status)
	return err
}

func (o *Order) UpdateNextStatus() {
	currentStatus := o.Statuses[len(o.Statuses)-1]
	o.Statuses = append(o.Statuses, currentStatus.Next())
	o.CurrentStatus = currentStatus.Next().Code
}

func (o *Order) UpdatePrevStatus() {
	currentStatus := o.Statuses[len(o.Statuses)-1]
	o.Statuses = append(o.Statuses, currentStatus.Prev())
	o.CurrentStatus = currentStatus.Prev().Code
}

type OrderItem struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	OrderID         *uint           `json:"order_id"`
	ProductID       *uint           `json:"product_id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Image           string          `json:"image"`
	SKU             string          `json:"sku"`
	Price           float64         `json:"price" gorm:"precision:18;scale:2"`
	Unit            string          `json:"unit"`
	Discount        float64         `json:"discount" gorm:"precision:18;scale:2"`
	DiscountedPrice float64         `json:"discounted_price" gorm:"precision:18;scale:2"` // DiscountPrice is the tax base after applying discount
	DiscountPercent float64         `json:"discount_percent" gorm:"precision:18;scale:2"`
	DiscountReason  string          `json:"discount_reason,omitempty"`
	Surcharge       float64         `json:"surcharge" gorm:"precision:18;scale:2"`
	SurchargeReason string          `json:"surcharge_reason,omitempty"`
	Comments        string          `json:"comments"`
	Course          string          `json:"course"`
	Hash            string          `json:"hash"`
	Modifiers       []OrderModifier `json:"modifiers"  gorm:"foreignKey:OrderItemID"`
	Tax             string          `json:"tax"`
	TaxPercentage   float64         `json:"tax_percentage"`
	TaxBase         float64         `json:"tax_base" gorm:"precision:18;scale:2"`
	TaxAmount       float64         `json:"tax_amount" gorm:"precision:18;scale:2"`
	CreatedAt       *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (oi *OrderItem) SetHash() {
	orderItemString := fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v", oi.ID, oi.OrderID, oi.Name, oi.Description, oi.Image, oi.Price, oi.Unit, oi.Discount, oi.DiscountReason, oi.Surcharge, oi.SurchargeReason, oi.Comments, oi.Course)
	for _, modifier := range oi.Modifiers {
		orderItemString += fmt.Sprintf("%v%v%v%v%v%v%v%v", modifier.ID, modifier.OrderItemID, modifier.Name, modifier.Description, modifier.Image, modifier.Price, modifier.Unit, modifier.Comments)
	}
	oi.Hash = fmt.Sprintf("%x", orderItemString)
}

func (oi *OrderItem) AddModifiers(modifier []OrderModifier) {
	oi.Modifiers = append(oi.Modifiers, modifier...)
}

func (oi *OrderItem) RemoveModifiers(modifiers []OrderModifier) {
	for _, modifier := range modifiers {
		for i, m := range oi.Modifiers {
			if m.ID == modifier.ID {
				oi.Modifiers = append(oi.Modifiers[:i], oi.Modifiers[i+1:]...)
				break
			}
		}
	}
}

type OrderModifier struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	OrderItemID     *uint           `json:"order_item_id"`
	OrderID         uint            `json:"order_id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Image           string          `json:"image"`
	Category        string          `json:"category"`
	ProductID       *uint           `json:"product_id"`
	SKU             string          `json:"sku"`
	Price           float64         `json:"price"  gorm:"precision:18;scale:2"`
	Discount        float64         `json:"discount" gorm:"precision:18;scale:2"`
	DiscountedPrice float64         `json:"discounted_price" gorm:"precision:18;scale:2"` // DiscountPrice is the tax base after applying discount
	DiscountPercent float64         `json:"discount_percent" gorm:"precision:18;scale:2"`
	DiscountReason  string          `json:"discount_reason,omitempty"`
	Surcharge       float64         `json:"surcharge" gorm:"precision:18;scale:2"`
	SurchargeReason string          `json:"surcharge_reason,omitempty"`
	Unit            string          `json:"unit"`
	Tax             string          `json:"tax"`
	TaxPercentage   float64         `json:"tax_percentage"`
	TaxAmount       float64         `json:"tax_amount" gorm:"precision:18;scale:2"`
	TaxBase         float64         `json:"tax_base" gorm:"precision:18;scale:2"`
	Comments        string          `json:"comments"`
	CreatedAt       *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt       *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type OrderType struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	ChannelID   *uint           `json:"channel_id"`
	StoreID     *uint           `json:"store_id"`
	BrandID     *uint           `json:"brand_id"`
	CreatedAt   *time.Time      `json:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty" swaggerignore:"true"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

type Attendee struct {
	ID        uint            `json:"id" swaggerignore:"true" gorm:"primaryKey"`
	AccountID uint            `json:"account_id"`
	OrderID   uint            `json:"order_id" swaggerignore:"true"`
	Name      string          `json:"name"`
	Role      string          `json:"role"`
	Action    OrderAction     `json:"order_action"`
	OrderStep OrderStep       `json:"order_step"`
	CreatedAt *time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt *time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`
}

type OrderStatus struct {
	ID        uint           `json:"id,omitempty" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"uniqueIndex:idx_order_status_code"`
	OrderID   *uint          `json:"order_id,omitempty" gorm:"uniqueIndex:idx_order_status_code"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" swaggerignore:"true"`
}

func (b *OrderStatus) BeforeCreate(tx *gorm.DB) (err error) {
	cols := []clause.Column{
		{
			Name: "code",
		},
		{
			Name: "order_id",
		},
	}

	tx.Statement.AddClause(clause.OnConflict{
		Columns:   cols,
		DoUpdates: clause.Assignments(map[string]interface{}{"updated_at": time.Now()}),
	})

	return nil
}

func (os *OrderStatus) Next() OrderStatus {
	switch os.Code {
	case OrderStatusCreated:
		return OrderStatus{
			Code: OrderStatusPaying,
		}
	case OrderStatusPaying:
		return OrderStatus{
			Code: OrderStatusClosed,
		}
	case OrderStatusClosed:
		return OrderStatus{
			Code: OrderStatusClosed,
		}
	default:
		return OrderStatus{
			Code: OrderStatusCreated,
		}
	}
}

func (os *OrderStatus) Prev() OrderStatus {
	switch os.Code {
	case OrderStatusCreated:
		return OrderStatus{
			Code: OrderStatusCreated,
		}
	case OrderStatusPaying:
		return OrderStatus{
			Code: OrderStatusCreated,
		}
	case OrderStatusClosed:
		return OrderStatus{
			Code: OrderStatusPaying,
		}
	default:
		return OrderStatus{
			Code: OrderStatusCreated,
		}
	}
}

type TipData struct {
	Percentage *float64 `json:"percentage"`
	Amount     *float64 `json:"value"`
}

func (t TipData) GetValueAndType() (string, float64) {
	if t.Percentage != nil && *t.Percentage > 0 {
		return "percentage", *t.Percentage / 100
	}

	if t.Amount != nil && *t.Amount > 0 {
		return "value", *t.Amount
	}

	return "", 0
}

type CloseInvoiceRequest struct {
	InvoiceID    string             `json:"invoice_id" swaggerignore:"true"`
	DocumentType string             `json:"document"`
	Payments     []*payment.Payment `json:"payments"`
	Observations string             `json:"observations"`
	attendee     *Attendee
}

func (c *CloseInvoiceRequest) GetTotalTips() float64 {
	total := 0.0
	for _, p := range c.Payments {
		total += p.Tip
	}
	return total
}
