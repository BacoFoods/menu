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
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type RequestInvoiceSplit struct {
	Invoices [][]uint `json:"invoices" binding:"required"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Get to handle a request to get an invoice
// @Tags Invoice
// @Summary To get an invoice
// @Description To get an invoice
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
// @Param paid query string false "paid" Enums(true,false)
// @Param store_id query string false "store id"
// @Param closed query string false "is closed" Enums(true,false)
// @Param days query string false "Days before"
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

	if storeID := c.Query("store_id"); storeID != "" {
		filter["store_id"] = storeID
	}

	if closed := c.Query("closed"); closed != "" {
		filter["closed"] = closed
	}

	if orderID := c.Query("order_id"); orderID != "" {
		filter["order_id"] = orderID
	}

	if paid := c.Query("paid"); paid != "" {
		filter["paid"] = paid
	}

	if days := c.Query("days"); days != "" {
		filter["days"] = days
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
// @Security ApiKeyAuth
// @Param id path string true "Invoice ID"
// @Param tip body RequestUpdateTip true "tip"
// @Success 200 {object} object{status=string,data=invoice.Invoice}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /invoice/{id} [post]
func (h *Handler) UpdateTip(c *gin.Context) {
	invoiceID := c.Param("id")

	var tipReq RequestUpdateTip
	if err := c.ShouldBindJSON(&tipReq); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateTip", err, tipReq)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	updatedInvoice, err := h.service.UpdateTip(tipReq.Value, tipReq.Type, invoiceID)
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

// Split to handle a request to split an invoice
// @Tags Invoice
// @Summary To split an invoice using items ids
// @Description To split an invoice using items ids
// @Param id path string true "invoice id"
// @Param invoices body RequestInvoiceSplit true "invoices for splitting"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=object{invoices=[]Invoice}}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /invoice/{id}/split [post]
func (h *Handler) Split(c *gin.Context) {
	invoiceID := c.Param("id")
	var body RequestInvoiceSplit
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "Split", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	invoices, err := h.service.Split(invoiceID, body.Invoices)
	if err != nil {
		shared.LogError("error separating invoice", LogHandler, "Split", err, invoices)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(invoices))
}

// Print to handle a request to print an invoice
// @Tags Invoice
// @Summary To print an invoice
// @Description To print an invoice
// @Param id path string true "invoice id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Invoice}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /invoice/{id}/print [get]
func (h *Handler) Print(c *gin.Context) {
	invoiceID := c.Param("id")

	printableInvoice, err := h.service.Print(invoiceID)
	if err != nil {
		shared.LogError("error printing invoice", LogHandler, "Print", err, invoiceID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorInvoicePrinting))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(printableInvoice))
}

// FindDiscountApplied to handle a request to find discount applied
// @Tags InvoiceApplied
// @Summary To find discount applied
// @Description To find discount applied
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=object{invoices=[]Invoice}}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /discount-applied [get]
func (h *Handler) FindDiscountApplied(c *gin.Context) {
	invoices, err := h.service.FindDiscountApplied()
	if err != nil {
		shared.LogError("error finding invoices", LogHandler, "FindDiscountApplied", err, invoices)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(invoices))
}

// RemoveDiscountApplied to handle a request to remove discount applied
// @Tags InvoiceApplied
// @Summary To remove discount applied
// @Description To remove discount applied
// @Param id path string true "invoice id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=object{invoices=[]DiscountApplied}}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /discount-applied/{id} [delete]
func (h *Handler) RemoveDiscountApplied(c *gin.Context) {
	invoiceAppliedID := c.Param("id")

	invoiceAppliedRemoved, err := h.service.RemoveDiscountApplied(invoiceAppliedID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(invoiceAppliedRemoved))
}

// DIAN Resolution

// FindResolution to handle a request to find resolutions
// @Tags Resolution
// @Summary To find resolutions
// @Description To find resolutions
// @Param storeID query string false "store id"
// @Param resolution query string false "resolution"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=object{resolutions=[]Resolution}}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /resolution [get]
func (h *Handler) FindResolution(c *gin.Context) {
	filter := make(map[string]any)

	if storeID := c.Query("store_id"); storeID != "" {
		filter["store_id"] = storeID
	}

	if resolution := c.Query("resolution"); resolution != "" {
		filter["resolution"] = resolution
	}

	resolutions, err := h.service.FindResolution(filter)
	if err != nil {
		shared.LogError("error finding resolutions", LogHandler, "FindResolution", err, resolutions)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(resolutions))
}

// CreateResolution to handle a request to create a resolution
// @Tags Resolution
// @Summary To create a resolution
// @Description To create a resolution
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param resolution body DTOResolution true "resolution"
// @Success 200 {object} object{status=string,data=Resolution}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /resolution [post]
func (h *Handler) CreateResolution(c *gin.Context) {
	var resolutionDTO DTOResolution
	if err := c.ShouldBindJSON(&resolutionDTO); err != nil {
		shared.LogError("error binding request body", LogHandler, "CreateResolution", err, resolutionDTO)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	resolution, err := resolutionDTO.ToResolution()
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	createdResolution, err := h.service.CreateResolution(resolution)
	if err != nil {
		shared.LogError("error creating resolution", LogHandler, "CreateResolution", err, createdResolution)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(createdResolution))
}

// UpdateResolution to handle a request to update a resolution
// @Tags Resolution
// @Summary To update a resolution
// @Description To update a resolution
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param resolution body Resolution true "resolution"
// @Success 200 {object} object{status=string,data=Resolution}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /resolution/{id} [patch]
func (h *Handler) UpdateResolution(c *gin.Context) {
	var resolution Resolution
	if err := c.ShouldBindJSON(&resolution); err != nil {
		shared.LogError("error binding request body", LogHandler, "UpdateResolution", err, resolution)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	updatedResolution, err := h.service.UpdateResolution(&resolution)
	if err != nil {
		shared.LogError("error updating resolution", LogHandler, "UpdateResolution", err, updatedResolution)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(updatedResolution))

}

// DeleteResolution to handle a request to delete a resolution
// @Tags Resolution
// @Summary To delete a resolution
// @Description To delete a resolution
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "resolution id"
// @Success 200 {object} object{status=string,data=Resolution}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /resolution/{id} [delete]
func (h *Handler) DeleteResolution(c *gin.Context) {
	resolutionID := c.Param("id")

	if err := h.service.DeleteResolution(resolutionID); err != nil {
		shared.LogError("error deleting resolution", LogHandler, "DeleteResolution", err, resolutionID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusNoContent, shared.SuccessResponse(nil))
}
