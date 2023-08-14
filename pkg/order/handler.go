package order

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const LogHandler string = "pkg/order/handler"

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
// @Param order body OrderTDP true "Order"
// @Success 200 {object} object{status=string,data=Order}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /order [post]
func (h *Handler) Create(c *gin.Context) {
	var body OrderTDP
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
