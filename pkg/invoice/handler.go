package invoice

import (
	"fmt"
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
	Type       string  `json:"type"`
	PaymentID  uint    `json:"payment_id"`
	Surcharges []uint  `json:"surcharge_id"`
	Tips       float64 `json:"tips"`
	Discounts  []uint  `json:"discount_id"`
}

type RequestUpdateTip struct {
	Tips float64 `json:"tips"`
}

func (r *RequestUpdateInvoice) ToInvoice() (*Invoice, error) {
	invoice, err := NewInvoiceBuilder().
		SetType(r.Type).
		SetPaymentID(r.PaymentID).
		SetTips(r.Tips).Build()
	if err != nil {
		return nil, err // Devuelve el error junto con nil
	}
	return invoice, nil
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
	fmt.Println("body", &body)

	updatedInvoice, err := body.ToInvoice()

	fmt.Println("updatedInvoice", updatedInvoice)

	if err != nil {
		shared.LogError("error ToInvoice", LogHandler, "ToInvoice", err, updatedInvoice)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvoiceUpdate))
		return
	}
	updatedInvoice, err = h.service.Update(updatedInvoice, body.Discounts, body.Surcharges) // TODO actualizar nombres a Discounts no DiscountID Actualización de la factura después de la creación

	if err != nil {
		shared.LogError("error updating invoice", LogHandler, "UpdateInvoice", err, updatedInvoice)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvoiceUpdate))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(updatedInvoice))
}

// UpdateTip to handle a request to update the tip of an invoice
// @Tags Invoice
// @Summary To update the tip of an invoice
// @Description To update the tip of an invoice
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param tip body RequestUpdateTip true "tip"
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /invoice/{id}/tip [post]
func (h *Handler) UpdateTip(c *gin.Context) {
	invoiceID := c.Param("id") // Obtener el ID del invoice desde el contexto

	var tipReq RequestUpdateTip
	if err := c.ShouldBindJSON(&tipReq); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateTip", err, tipReq)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	// Obtener el invoice existente por su ID
	existingInvoice, err := h.service.Get(invoiceID)
	if err != nil {
		shared.LogError("error getting invoice", LogHandler, "UpdateTip", err, existingInvoice)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingInvoice))
		return
	}

	// Verificar si las propinas son nulas o negativas
	if tipReq.Tips < 0 {
		shared.LogError("error invalid tip amount", LogHandler, "UpdateTip", err, tipReq.Tips)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvalidTipAmount))
		return
	}

	// Calcular el subtotal
	subtotal := existingInvoice.SubTotal

	// Verificar si los tips son un porcentaje o un valor nominal
	if tipReq.Tips <= 1.0 { // Si es menor o igual a 1, se considera un porcentaje
		// Verificar si el porcentaje excede el 10% del subtotal
		if tipReq.Tips > 0.1*subtotal {
			shared.LogError("error tip percentage exceeds limit", LogHandler, "UpdateTip", err, tipReq.Tips)
			c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTipPercentageExceedsLimit))
			return
		}
		existingInvoice.Tips = tipReq.Tips * subtotal // Calcular las propinas como un porcentaje del subtotal
	} else { // Si es mayor que 1, se considera un valor nominal y se suma directamente
		existingInvoice.Tips += tipReq.Tips
	}

	// Guardar el invoice actualizado en la base de datos
	updatedInvoice, err := h.service.Update(existingInvoice, nil, nil) // No se actualizan descuentos ni recargos
	if err != nil {
		shared.LogError("error updating invoice", LogHandler, "UpdateTip", err, updatedInvoice)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvoiceUpdate))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(updatedInvoice))
}
