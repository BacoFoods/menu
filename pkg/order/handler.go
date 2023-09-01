package order

import (
	"fmt"
	"net/http"
	"strconv"

	invoices "github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const LogHandler string = "pkg/order/handler"

type RequestUpdateOrderSeats struct {
	Seats int `json:"seats" binding:"required"`
}

type RequestUpdateOrderStatus struct {
	Status string `json:"status" binding:"required" enum:"ordering,preparing,delivered,invoicing,canceled,completed" example:"ordering,preparing,delivered,invoicing,canceled,completed"`
}

type RequestUpdateOrderProduct struct {
	Price    float64 `json:"price"`
	Unit     string  `json:"unit"`
	Quantity int     `json:"quantity"`
	Comments string  `json:"comments"`
	Course   string  `json:"course"`
}

type RequestAddProducts struct {
	Items []OrderItemDTO `json:"items" binding:"required"`
}

type RequestModifiers struct {
	Modifiers []OrderModifierDTO `json:"modifiers" binding:"required"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Order

// Create to handle a request to create an order
// @Tags Order
// @Summary To create an order
// @Description To create an order
// @Accept json
// @Produce json
// @Param order body OrderDTO true "Order"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order [post]
func (h *Handler) Create(c *gin.Context) {
	var body OrderDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "Create", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	order := body.ToOrder()
	orderDB, err := h.service.Create(&order)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(orderDB))
}

// UpdateTable to handle a request to update the table of an order
// @Tags Order
// @Summary To update the table of an order
// @Description To update the table of an order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param table path string true "Table"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order/{id}/table/{table} [patch]
func (h *Handler) UpdateTable(c *gin.Context) {
	storeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		shared.LogWarn("error parsing store id", LogHandler, "UpdateTable", err, storeID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequestStoreID))
		return
	}

	tableID, err := strconv.ParseUint(c.Param("table"), 10, 64)
	if err != nil {
		shared.LogWarn("error parsing table id", LogHandler, "UpdateTable", err, tableID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequestTableID))
		return
	}

	order, err := h.service.UpdateTable(storeID, tableID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// Get to handle a request to get an order
// @Tags Order
// @Summary To get an order
// @Description To get an order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	orderID := c.Param("id")

	order, err := h.service.Get(orderID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// Find to handle a request to find orders
// @Tags Order
// @Summary To find orders
// @Description To find orders
// @Accept json
// @Produce json
// @Param store query string false "Store ID"
// @Param table query string false "Table ID"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order [get]
func (h *Handler) Find(c *gin.Context) {
	filters := make(map[string]interface{})

	storeID := c.Query("store")
	if storeID != "" {
		filters["store_id"] = storeID
	}

	tableID := c.Query("table")
	if tableID != "" {
		filters["table_id"] = tableID
	}

	orders, err := h.service.Find(filters)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(orders))
}

// UpdateSeats to handle a request to update the seats of an order
// @Tags Order
// @Summary To update the seats of an order
// @Description To update the seats of an order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param seats body RequestUpdateOrderSeats true "Seats"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order/{id}/seats [patch]
func (h *Handler) UpdateSeats(c *gin.Context) {
	var body RequestUpdateOrderSeats
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateSeats", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	if body.Seats < 0 {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequestOrderSeats))
		return
	}

	orderID := c.Param("id")

	order, err := h.service.UpdateSeats(orderID, body.Seats)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// AddProducts to handle a request to add a products to an order
// @Tags Order
// @Summary To add products to an order
// @Description To add products to an order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param product body RequestAddProducts true "Add Products"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order/{id}/add/products [patch]
func (h *Handler) AddProducts(c *gin.Context) {
	orderID := c.Param("id")

	var body RequestAddProducts
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "AddProduct", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	items := make([]OrderItem, 0)
	for _, item := range body.Items {
		items = append(items, item.ToOrderItem())
	}

	order, err := h.service.AddProducts(orderID, items)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// RemoveProduct to handle a request to remove a product from an order
// @Tags Order
// @Summary To remove a product from an order
// @Description To remove a product from an order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param productID path string true "Product ID"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order/{id}/remove/product [patch]
func (h *Handler) RemoveProduct(c *gin.Context) {
	orderID := c.Param("id")
	productID := c.Param("productID")

	order, err := h.service.RemoveProduct(orderID, productID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// UpdateProduct to handle a request to update a product from an order
// @Tags Order
// @Summary To update a product from an order
// @Description To update a product from an order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param productID path string true "Product ID"
// @Param product body RequestUpdateOrderProduct true "product"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order/{id}/update/product [patch]
func (h *Handler) UpdateProduct(c *gin.Context) {
	var body RequestUpdateOrderProduct
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateProduct", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequestOrderID))
		return
	}

	productID, err := strconv.ParseUint(c.Param("productID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequestProductID))
		return
	}

	uOrderID := uint(orderID)
	uProductID := uint(productID)
	updatedProduct := &OrderItem{
		OrderID:   &uOrderID,
		ProductID: &uProductID,
		Price:     body.Price,
		Unit:      body.Unit,
		Comments:  body.Comments,
		Course:    body.Course,
	}

	order, err := h.service.UpdateProduct(updatedProduct)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// UpdateStatus to handle a request to update the status of an order
// @Tags Order
// @Summary To update the status of an order
// @Description To update the status of an order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body RequestUpdateOrderStatus true "Status"
// @Success 200 {object} object{status=string,data=Order}
// @Router /order/{id}/status [patch]
func (h *Handler) UpdateStatus(c *gin.Context) {
	orderID := c.Param("id")

	var body RequestUpdateOrderStatus
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateStatus", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	order, err := h.service.UpdateStatus(orderID, body.Status)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// ReleaseTable to handle a request to release an order's table
// @Tags Order
// @Summary To release an order's table
// @Description To release an order's table
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /order/{id}/release-table [post]
func (h *Handler) ReleaseTable(c *gin.Context) {
	orderID := c.Param("id")

	order, err := h.service.ReleaseTable(orderID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// Order Items

// AddModifiers to handle a request to add a modifiers to a product's order
// @Tags Order
// @Summary To add modifiers to a product's order
// @Description To add modifiers to a product's order
// @Accept json
// @Produce json
// @Param id path string true "OrderItemID"
// @Param modifier body RequestModifiers true "Add Modifiers"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order-item/{id}/add/modifiers [patch]
func (h *Handler) AddModifiers(c *gin.Context) {
	orderItemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequestOrderItemID))
		return
	}

	var body RequestModifiers
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "AddModifiers", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	modifiers := make([]OrderModifier, 0)
	for _, modifier := range body.Modifiers {
		modifiers = append(modifiers, modifier.ToOrderModifier())
	}

	if len(modifiers) == 0 {
		shared.LogError("no modifiers provided", LogHandler, "AddModifiers", nil, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	order, err := h.service.AddModifiers(uint(orderItemID), modifiers)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// RemoveModifiers to handle a request to remove a modifiers from a product's order
// @Tags Order
// @Summary To remove a modifiers from a product's order
// @Description To remove a modifiers from a product's order
// @Accept json
// @Produce json
// @Param id path string true "OrderItemID"
// @Param modifier body RequestModifiers true "Remove Modifiers"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order-item/{id}/remove/modifiers [patch]
func (h *Handler) RemoveModifiers(c *gin.Context) {
	orderItemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequestOrderItemID))
		return
	}

	var body RequestModifiers
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "RemoveModifiers", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	modifiers := make([]OrderModifier, 0)
	for _, modifier := range body.Modifiers {
		modifiers = append(modifiers, modifier.ToOrderModifier())
	}

	if len(modifiers) == 0 {
		shared.LogError("no modifiers provided", LogHandler, "RemoveModifiers", nil, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	order, err := h.service.RemoveModifiers(uint(orderItemID), modifiers)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// Oder Types

// CreateOrderType to handle a request to create an order type
// @Tags OrderType
// @Summary To create an order type
// @Description To create an order type
// @Accept json
// @Produce json
// @Param orderType body OrderType true "Order Type"
// @Success 200 {object} object{status=string,data=OrderType}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order-type [post]
func (h *Handler) CreateOrderType(c *gin.Context) {
	var body OrderType
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "CreateOrderType", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	orderType, err := h.service.CreateOrderType(&body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(orderType))
}

// FindOrderType to handle a request to find order types
// @Tags OrderType
// @Summary To find order types
// @Description To find order types
// @Accept json
// @Produce json
// @Param channelID path string false "Channel ID"
// @Param storeID path string false "Store ID"
// @Param brandID path string false "Brand ID"
// @Param name path string false "Name"
// @Success 200 {object} object{status=string,data=OrderType}
// @Router /order-type [get]
func (h *Handler) FindOrderType(c *gin.Context) {
	filters := make(map[string]interface{})

	channelID := c.Query("channelID")
	if channelID != "" {
		filters["channel_id"] = channelID
	}

	storeID := c.Query("storeID")
	if storeID != "" {
		filters["store_id"] = storeID
	}

	brandID := c.Query("brandID")
	if brandID != "" {
		filters["brand_id"] = brandID
	}

	name := c.Query("name")
	if name != "" {
		filters["name"] = name
	}

	orderTypes, err := h.service.FindOrderType(filters)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(orderTypes))
}

// GetOrderType to handle a request to get an order type
// @Tags OrderType
// @Summary To get an order type
// @Description To get an order type
// @Accept json
// @Produce json
// @Param id path string true "Order Type ID"
// @Success 200 {object} object{status=string,data=OrderType}
// @Router /order-type/{id} [get]
func (h *Handler) GetOrderType(c *gin.Context) {
	orderTypeID := c.Param("id")

	orderType, err := h.service.GetOrderType(orderTypeID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(orderType))
}

// UpdateOrderType to handle a request to update an order type
// @Tags OrderType
// @Summary To update an order type
// @Description To update an order type
// @Accept json
// @Produce json
// @Param id path string true "Order Type ID"
// @Param orderType body OrderType true "Order Type"
// @Success 200 {object} object{status=string,data=OrderType}
// @Router /order-type/{id} [patch]
func (h *Handler) UpdateOrderType(c *gin.Context) {
	orderTypeID := c.Param("id")

	var body OrderType
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateOrderType", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	orderType, err := h.service.UpdateOrderType(orderTypeID, &body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(orderType))
}

// DeleteOrderType to handle a request to delete an order type
// @Tags OrderType
// @Summary To delete an order type
// @Description To delete an order type
// @Accept json
// @Produce json
// @Param id path string true "Order Type ID"
// @Success 200 {object} object{status=string,data=OrderType}
// @Router /order-type/{id} [delete]
func (h *Handler) DeleteOrderType(c *gin.Context) {
	orderTypeID := c.Param("id")

	err := h.service.DeleteOrderType(orderTypeID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(fmt.Sprintf("Order type with id %s has been deleted", orderTypeID)))
}

// Invoice

// CreateInvoice to handle a request to create an invoice
// @Tags Invoice
// @Summary To create an invoice
// @Description To create an invoice
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Router /invoice [post]
func (h *Handler) CreateInvoice(c *gin.Context) {
	orderID := c.Param("id")

	invoice, err := h.service.CreateInvoice(orderID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(invoices.ErrorInvoiceCreation))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(invoice))
}
