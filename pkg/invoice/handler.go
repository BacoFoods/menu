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

type RequestUpdateTip struct {
	Tips float64 `json:"tips"`
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
// @Failure 401 {object} shared.Response
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

// Find to handle a request to find invoices
// @Tags Invoice
// @Summary To find invoices
// @Description To find invoices
// @Param order_id query string false "order id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=object{invoices=[]Invoice}}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /invoice [get]
func (h *Handler) Find(c *gin.Context) {
	filter := make(map[string]any)
	orderID := c.Query("order_id")
	if orderID != "" {
		filter["order_id"] = orderID
	}

	invoices, err := h.service.Find(filter)
	if err != nil {
		shared.LogError("error finding invoices", LogHandler, "Find", err, invoices)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvoiceFinding))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(invoices))
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
// @Router /invoice/{id} [post]
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

	// Actualizar el campo 'tips' del invoice existente
	existingInvoice.Tips = tipReq.Tips

	// Guardar el invoice actualizado en la base de datos
	updatedInvoice, err := h.service.UpdateTip(existingInvoice) // No se actualizan descuentos ni recargos
	if err != nil {
		shared.LogError("error updating invoice", LogHandler, "UpdateTip", err, updatedInvoice)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvoiceUpdate))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(updatedInvoice))
}

// AddClient to handle a request to add a client to an invoice
// @Tags Invoice
// @Summary To add a client to an invoice
// @Description To add a client to an invoice
// @Param id path string true "invoice id"
// @Param clientID path string true "client id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Invoice}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /invoice/{id}/client/{clientID}/add [post]
func (h *Handler) AddClient(c *gin.Context) {
	invoiceID := c.Param("id")
	clientID := c.Param("clientID")

	invoice, err := h.service.AddClient(invoiceID, clientID)
	if err != nil {
		shared.LogError("error adding client to invoice", LogHandler, "AddClient", err, invoice)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvoiceAddingClient))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(invoice))
}

// RemoveClient to handle a request to remove a client from an invoice
// @Tags Invoice
// @Summary To remove a client from an invoice
// @Description To remove a client from an invoice
// @Param id path string true "invoice id"
// @Param clientID path string true "client id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Invoice}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /invoice/{id}/client/{clientID}/remove [post]
func (h *Handler) RemoveClient(c *gin.Context) {
	invoiceID := c.Param("id")
	clientID := c.Param("clientID")

	invoice, err := h.service.RemoveClient(invoiceID, clientID)
	if err != nil {
		shared.LogError("error removing client from invoice", LogHandler, "RemoveClient", err, invoice)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvoiceRemovingClient))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(invoice))
}
