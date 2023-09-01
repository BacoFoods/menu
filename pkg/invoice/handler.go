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

type RequestUpdateInvoice struct {
	Type         string `json:"type"`
	PaymentID    uint   `json:"payment_id"`
	SurchargeID  uint   `json:"surcharge_id"`
	Tips         float64 `json:"tips"`
	DiscountID   uint    `json:"discount_id"`
}


func (r *RequestUpdateInvoice) ToInvoice() *Invoice {
	discount := Discount{
		ID: r.DiscountID,
	}
	surcharge := Surcharge{
		ID: r.SurchargeID,
	}

	paymentID := r.PaymentID
	var paymentIDPtr *uint

	if paymentID != 0 {
		paymentIDPtr = &paymentID
	}

	return &Invoice{
		Type:      r.Type,
		PaymentID: paymentIDPtr,
		Discounts: []Discount{discount},
		Surcharges: []Surcharge{surcharge},
	}
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Get to handle a request to get an invoice
// @Tags Invoice
// @Summary To get a invoice
// @Description To get a invoice
// @Param id path string true "Invoice ID"
// @Accept json
// @Produce json
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

// UpdateInvoice to handle a request to update an invoice
// @Tags Invoice
// @Summary To update an invoice
// @Description To update an invoice
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param invoice body RequestUpdateInvoice true "invoice"
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /invoice/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	var body RequestUpdateInvoice
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateInvoice", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	updatedInvoice, err := h.service.Update(body.ToInvoice())

	if err != nil {
		shared.LogError("error updating invoice", LogHandler, "UpdateInvoice", err, updatedInvoice)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvoiceUpdate))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(updatedInvoice))
}