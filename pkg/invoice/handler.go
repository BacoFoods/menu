package invoice

import (
	"net/http"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const LogHandler = "pkg/invoice/handler"

type CreateInvoiceRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Get to handle a request to get a invoice
// @Tags Invoice
// @Summary To get a invoice
// @Description To get a invoice
// @Param id path string true "invoice id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /invoice/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	invoiceID := c.Param("id")

	invoices, err := h.service.Get(invoiceID)
	if err != nil {
		shared.LogError("error getting invoice", LogHandler, "Get", err, invoices)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingInvoice))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(invoices))
}
