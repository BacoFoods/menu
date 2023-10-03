package payment

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	LogHandler string = "pkg/payment/handler"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Get to handle a request to get a payment by id
// @Tags Payment
// @Summary To get a payment by id
// @Description To get a payment by id
// @Param id path string true "payment id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Payment}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	paymentID := ctx.Param("id")

	payment, err := h.service.Get(paymentID)
	if err != nil {
		shared.LogError("error getting payment", LogHandler, "Get", err, paymentID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentGetting))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(payment))
}

// Find to handle a request to find all payment
// @Tags Payment
// @Summary To find payment
// @Description To find payment
// @Param code query string false "payment code"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Payment}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment [get]
func (h *Handler) Find(ctx *gin.Context) {
	var filter map[string]any

	if code := ctx.Query("code"); code != "" {
		filter["code"] = code
	}

	payments, err := h.service.Find(filter)
	if err != nil {
		shared.LogError("error finding payment", LogHandler, "Find", err)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentFinding))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(payments))
}

// Create to handle a request to create payment
// @Tags Payment
// @Summary To create payment
// @Description To create payment
// @Param payment body Payment true "payment body"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Payment}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment [post]
func (h *Handler) Create(ctx *gin.Context) {
	var body Payment
	if err := ctx.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding json", LogHandler, "Create", err)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentCreating))
		return
	}

	payment, err := h.service.Create(&body)
	if err != nil {
		shared.LogError("error creating payment", LogHandler, "Create", err, body)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentCreating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(payment))
}

// Update to handle a request to update payment
// @Tags Payment
// @Summary To update payment
// @Description To update payment
// @Param payment body Payment true "payment body"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Payment}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment [patch]
func (h *Handler) Update(ctx *gin.Context) {
	var body Payment
	if err := ctx.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding json", LogHandler, "Update", err, body)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentUpdating))
		return
	}

	payment, err := h.service.Update(&body)
	if err != nil {
		shared.LogError("error updating payment", LogHandler, "Update", err, body)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentUpdating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(payment))
}

// Delete to handle a request to delete payment
// @Tags Payment
// @Summary To delete payment
// @Description To delete payment
// @Param id path string true "payment id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Payment}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	paymentID := ctx.Param("id")

	payment, err := h.service.Delete(paymentID)
	if err != nil {
		shared.LogError("error deleting payment", LogHandler, "Delete", err, paymentID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentDeleting))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(payment))
}

// FindPaymentMethod to handle a request to find all payment method
// @Tags PaymentMethod
// @Summary To find payment method
// @Description To find payment method
// @Accept json
// @Produce json
// @Param brand_id query string false "brand id"
// @Param store_id query string false "store id"
// @Param channel_id query string false "channel id"
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]PaymentMethod}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment-method [get]
func (h *Handler) FindPaymentMethod(ctx *gin.Context) {
	var filter map[string]any

	if brandID := ctx.Query("brand_id"); brandID != "" {
		filter["brand_id"] = brandID
	}

	if storeID := ctx.Query("store_id"); storeID != "" {
		filter["store_id"] = storeID
	}

	if channelID := ctx.Query("channel_id"); channelID != "" {
		filter["channel_id"] = channelID
	}

	paymentMethods, err := h.service.FindPaymentMethods(filter)
	if err != nil {
		shared.LogError("error finding payment method", LogHandler, "FindPaymentMethod", err)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentMethodFinding))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(paymentMethods))
}

// GetPaymentMethod to handle a request to get a payment method by id
// @Tags PaymentMethod
// @Summary To get a payment method by id
// @Description To get a payment method by id
// @Param id path string true "payment method id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=PaymentMethod}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment-method/{id} [get]
func (h *Handler) GetPaymentMethod(ctx *gin.Context) {
	paymentMethodID := ctx.Param("id")

	paymentMethod, err := h.service.GetPaymentMethod(paymentMethodID)
	if err != nil {
		shared.LogError("error getting payment method", LogHandler, "GetPaymentMethod", err, paymentMethodID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentMethodFinding))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(paymentMethod))
}

// CreatePaymentMethod to handle a request to create payment method
// @Tags PaymentMethod
// @Summary To create payment method
// @Description To create payment method
// @Param payment_method body PaymentMethod true "payment method body"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=PaymentMethod}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment-method [post]
func (h *Handler) CreatePaymentMethod(ctx *gin.Context) {
	var body PaymentMethod
	if err := ctx.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding json", LogHandler, "CreatePaymentMethod", err, body)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentCreating))
		return
	}

	paymentMethod, err := h.service.CreatePaymentMethod(&body)
	if err != nil {
		shared.LogError("error creating payment method", LogHandler, "CreatePaymentMethod", err, body)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentCreating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(paymentMethod))
}

// UpdatePaymentMethod to handle a request to update payment method
// @Tags PaymentMethod
// @Summary To update payment method
// @Description To update payment method
// @Param payment_method body PaymentMethod true "payment method body"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=PaymentMethod}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment-method [patch]
func (h *Handler) UpdatePaymentMethod(ctx *gin.Context) {
	var body PaymentMethod
	if err := ctx.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding json", LogHandler, "UpdatePaymentMethod", err, body)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentUpdating))
		return
	}

	paymentMethod, err := h.service.UpdatePaymentMethod(&body)
	if err != nil {
		shared.LogError("error updating payment method", LogHandler, "UpdatePaymentMethod", err, body)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentUpdating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(paymentMethod))

}

// DeletePaymentMethod to handle a request to delete payment method
// @Tags PaymentMethod
// @Summary To delete payment method
// @Description To delete payment method
// @Param id path string true "payment method id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=PaymentMethod}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /payment-method/{id} [delete]
func (h *Handler) DeletePaymentMethod(ctx *gin.Context) {
	paymentMethodID := ctx.Param("id")

	paymentMethod, err := h.service.DeletePaymentMethod(paymentMethodID)
	if err != nil {
		shared.LogError("error deleting payment method", LogHandler, "DeletePaymentMethod", err, paymentMethodID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorPaymentDeleting))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(paymentMethod))
}
