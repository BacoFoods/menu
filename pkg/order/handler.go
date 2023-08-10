package order

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/order/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

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
