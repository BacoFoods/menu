package order

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BacoFoods/menu/pkg/client"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const LogHandler string = "pkg/order/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Order

func (h *Handler) ctxAttendee(ctx context.Context) *Attendee {
	username := ""
	if ctx.Value("account_name") != nil {
		username = ctx.Value("account_name").(string)
	}

	role := ""
	if ctx.Value("role") != nil {
		role = ctx.Value("role").(string)
	}
	accountID := ""
	if ctx.Value("account_id") != nil {
		accountID = ctx.Value("account_id").(string)
	}

	id, err := strconv.Atoi(accountID)
	if err != nil {
		shared.LogError("error parsing account id", LogHandler, "ctxAttendee", err, accountID)
		return nil
	}

	return &Attendee{
		AccountID: uint(id),
		Name:      username,
		Role:      role,
	}
}

// Create to handle a request to create an order
// @Tags Order
// @Summary To create an order
// @Description To create an order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param order body OrderDTO true "Order"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /order [post]
func (h *Handler) Create(c *gin.Context) {
	var body OrderDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "Create", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	idempotencyKey := c.GetHeader("X-Idempotence-Key")

	order := body.ToOrder()
	orderDB, err := h.service.Create(idempotencyKey, &order, c)
	if err == errDuplicatedOrder {
		c.JSON(http.StatusLocked, shared.ErrorResponse(err.Error()))
		return
	}

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(orderDB))
}

// Update to handle a request to update an order
// @Tags Order
// @Summary to update an order
// @Description to update an order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param order body OrderDTO true "Order"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /order [patch]
func (h *Handler) Update(c *gin.Context) {
	var body OrderDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "Update", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
	}

	order := body.ToOrder()
	orderUpdated, err := h.service.Update(&order)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(orderUpdated))
}

// CreatePublic to handle a request to create an order
// @Tags Order
// @Summary To create an order
// @Description To create an order
// @Accept json
// @Produce json
// @Param order body OrderDTO true "Order"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /public/order [post]
func (h *Handler) CreatePublic(c *gin.Context) {
	var body OrderDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "Create", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	order := body.ToOrder()
	orderDB, err := h.service.Create("", &order, c)
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
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param table path string true "Table"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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

// GetPublic to handle a request to get an order
// @Tags Order
// @Summary To get an order
// @Description To get an order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /public/order/{id} [get]
func (h *Handler) GetPublic(c *gin.Context) {
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
// @Security ApiKeyAuth
// @Param store query string false "Store ID"
// @Param table query string false "Table ID"
// @Param status query string false "Status" Enums(created, paying, closed)
// @Param days query string false "Days before"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /order [get]
func (h *Handler) Find(c *gin.Context) {
	filters := make(map[string]interface{})

	if storeID := c.Query("store"); storeID != "" {
		filters["store_id"] = storeID
	}

	if tableID := c.Query("table"); tableID != "" {
		filters["table_id"] = tableID
	}

	if status := c.Query("status"); status != "" {
		filters["current_status"] = strings.Split(status, ",")
	}

	if days := c.Query("days"); days != "" {
		filters["days"] = days
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
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param seats body RequestUpdateOrderSeats true "Seats"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param product body RequestAddProducts true "Add Products"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /order/{id}/add/products [patch]
func (h *Handler) AddProducts(c *gin.Context) {
	orderID := c.Param("id")

	var body RequestAddProducts
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "AddProduct", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	idempoKey := c.GetHeader("X-Idempotence-Key")

	items := make([]OrderItem, 0)
	for _, item := range body.Items {
		for j := 0; j < item.Quantity; j++ {
			items = append(items, item.ToOrderItem())
		}
	}

	order, err := h.service.AddProducts(idempoKey, orderID, items)
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
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param productID path string true "Product ID"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param productID path string true "Product ID"
// @Param product body RequestUpdateOrderProduct true "product"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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

// UpdateComments to handle a request to update an order's comments
// @Tags Order
// @Summary To update an order's comments
// @Description To update an order's comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param comments body RequestUpdateOrderComments true "Comments"
// @Success 200 {object} object{status=string,data=Order}
// @Router /order/{id}/update/comments [patch]
func (h *Handler) UpdateComments(c *gin.Context) {
	orderID := c.Param("id")

	var body RequestUpdateOrderComments
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateComments", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	order, err := h.service.UpdateComments(orderID, body.Comments)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOrderUpdatingComments))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// UpdateClientName to handle a request to update an order's client name
// @Tags Order
// @Summary To update an order's client name
// @Description To update an order's client name
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param clientName body RequestUpdateOrderClientName true "Client Name"
// @Success 200 {object} object{status=string,data=Order}
// @Router /order/{id}/update/client-name [patch]
func (h *Handler) UpdateClientName(c *gin.Context) {
	orderID := c.Param("id")

	var body RequestUpdateOrderClientName
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateClientName", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	order, err := h.service.UpdateClientName(orderID, body.ClientName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOrderUpdatingClientName))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(order))
}

// UpdateStatus to handle a request to update an order's status
// @Tags Order
// @Summary To update an order's status
// @Description To update an order's status
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param status body RequestUpdateOrderStatus true "Status"
// @Success 200 {object} object{status=string,data=Order}
// @Router /order/{id}/update/status [patch]
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
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOrderUpdatingStatus))
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
// @Security ApiKeyAuth
// @Param id path string true "OrderItemID"
// @Param modifier body RequestModifiers true "Add Modifiers"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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
// @Security ApiKeyAuth
// @Param id path string true "OrderItemID"
// @Param modifier body RequestModifiers true "Remove Modifiers"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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

// OrderItemUpdate to handle a request to update an order item
// @Tags Order
// @Summary To update an order item
// @Description To update an order item
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "OrderItemID"
// @Param orderItem body RequestUpdateOrderItem true "order item"
// @Success 200 {object} object{status=string,data=Order}
// @Router /order-item/{id}/update [patch]
func (h *Handler) OrderItemUpdate(c *gin.Context) {
	orderItemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		shared.LogWarn("error parsing order item id", LogHandler, "OrderItemUpdate", err, orderItemID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequestOrderItemID))
	}

	var body RequestUpdateOrderItem
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateOrderItem", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	item := body.ToOrderItem()
	item.ID = uint(orderItemID)

	orderItem, err := h.service.OrderItemUpdateCourse(&item)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOrderItemUpdateCourse))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(orderItem))
}

// Oder Types

// CreateOrderType to handle a request to create an order type
// @Tags OrderType
// @Summary To create an order type
// @Description To create an order type
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param orderType body OrderType true "Order Type"
// @Success 200 {object} object{status=string,data=OrderType}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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

type CreateInvoiceDocumentRequest struct {
	// Possible options are: POS, FEUnidentified, FEIdentified
	DocumentType string `json:"document_type"`

	// Client data for the invoice
	DocumentData *client.Client `json:"document_data"`
}

type CreateInvoiceRequest struct {
	RequestCalculateInvoice

	CreateInvoiceDocumentRequest

	RequestInvoicePaymentMethod

	orderId  string
	attendee *Attendee
}

// CreateInvoice to handle a request to create an invoice
// @Tags Order
// @Summary To create an invoice
// @Description To create an invoice
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param body body CreateInvoiceRequest true "request body"
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Router /order/{id}/invoice [post]
func (h *Handler) CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.LogError("error binding request body", LogHandler, "CreateInvoice", err, req)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	req.orderId = c.Param("id")
	req.attendee = h.ctxAttendee(c)

	// TODO: improve payment method
	req.PaymentMethodID = 1

	invoiceDB, err := h.service.CreateInvoice(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(invoiceDB))
}

// CalculateInvoice to handle a request to calculate an invoice
// @Tags Order
// @Summary To calculate an invoice
// @Description To calculate an invoice
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Param body body RequestCalculateInvoice true "request body"
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Router /order/{id}/invoice/calculate [post]
func (h *Handler) CalculateInvoice(c *gin.Context) {
	orderID := c.Param("id")
	var req RequestCalculateInvoice
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.LogError("error binding request body", LogHandler, "CalculateInvoice", err, req)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	newInvoice, err := h.service.CalculateInvoice(orderID, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOrderInvoiceCalculation))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(newInvoice))
}

// PublicCalculateInvoice to handle a request to calculate an invoice
// @Tags Order
// @Summary Public endpoint to calculate an invoice
// @Description Public endpoint to calculate an invoice
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Router /public/order/{id}/invoice/calculate [get]
func (h *Handler) PublicCalculateInvoice(c *gin.Context) {
	orderID := c.Param("id")

	n, o, err := h.service.CalculateInvoiceOIT(orderID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOrderInvoiceCalculation))
		return
	}

	res := struct {
		OldInvoice *invoice.Invoice `json:"old_invoice"`
		NewInvoice *invoice.Invoice `json:"new_invoice"`
	}{o, n}

	c.JSON(http.StatusOK, shared.SuccessResponse(res))
}

type CheckoutRequest struct {
	Tip        float64 `json:"tip"`
	CustomerID *string `json:"customer_id"`
}

// PublicCheckout to handle the checkout process of an order. Public for OIT.
// @Tags Order
// @Summary To check out an order.
// @Description To check out an order.
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param body body CheckoutRequest true "Checkout parameters"
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Router /public/order/{id}/checkout [post]
func (h *Handler) PublicCheckout(c *gin.Context) {
	// TODO: not public, but using user bearer token from ecom
	var checkout CheckoutRequest
	orderID := c.Param("id")

	if err := c.ShouldBindJSON(&checkout); err != nil {
		shared.LogError("error binding request body", LogHandler, "PublicCheckout", err, checkout)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	newInvoice, err := h.service.Checkout(orderID, checkout)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOrderInvoiceCalculation))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(newInvoice))
}

// CloseInvoice to handle a request to close an invoice
// @Tags Invoice
// @Summary To close an invoice
// @Description To close an invoice
// @Param id path string true "invoice id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param invoice body CloseInvoiceRequest true "invoice"
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /invoice/{id}/close [post]
func (h *Handler) CloseInvoice(c *gin.Context) {
	var req CloseInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.LogError("error binding request body", LogHandler, "CloseInvoice", err, req)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	req.InvoiceID = c.Param("id")
	att := h.ctxAttendee(c)
	req.attendee = att

	newInvoice, err := h.service.CloseInvoice(req)
	if err != nil {
		shared.LogError("error closing invoice", LogHandler, "CloseInvoice", err, newInvoice)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(newInvoice))
}
