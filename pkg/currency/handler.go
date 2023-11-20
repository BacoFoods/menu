package currency

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/currency/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create to handle a request for create a currency
// @Tags Currency
// @Summary To create a currency
// @Description To create a currency
// @Param currency body Currency true "currency request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Currency}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /currency [post]
func (h *Handler) Create(ctx *gin.Context) {
	var requestBody Currency
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCurrencyBadRequest))
		return
	}

	currency, err := h.service.Create(&requestBody)
	if err != nil {
		shared.LogError("error creating currency", LogHandler, "Create", err, currency)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCurrencyCreation))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(currency))
}

// Find to handle a request for find all currency
// @Tags Currency
// @Summary To find currency
// @Description To find currency
// @Param code query string false "currency code"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Currency}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /currency [get]
func (h *Handler) Find(ctx *gin.Context) {
	query := make(map[string]string)
	code := ctx.Query("code")
	if code != "" {
		query["code"] = code
	}
	currency, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error getting all currency", LogHandler, "Find", err)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCurrencyGetting))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(currency))
}

// Get to handle a request for find a currency by id
// @Tags Currency
// @Summary To find currency by id
// @Description To find currency by id
// @Param id path string true "currency id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Currency}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /currency/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	currencyID := ctx.Param("id")
	currency, err := h.service.Get(currencyID)
	if err != nil {
		shared.LogError("error getting currency", LogHandler, "Get", err, currencyID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCurrencyGetting))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(currency))
}

// Update to handle a request for update a currency
// @Tags Currency
// @Summary To update currency
// @Description To update currency plan
// @Param currency body Currency true "currency to update"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Currency}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /currency [patch]
func (h *Handler) Update(ctx *gin.Context) {
	var requestBody Currency
	if err := ctx.BindJSON(&requestBody); err != nil {
		shared.LogError("error getting currency request body", LogHandler, "Update", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCurrencyBadRequest))
		return
	}
	currency, err := h.service.Update(requestBody)
	if err != nil {
		shared.LogError("error updating currency", LogHandler, "Update", err, currency)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCurrencyUpdating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(currency))
}

// Delete to handle a request for delete a currency
// @Tags Currency
// @Summary To delete a currency
// @Description To delete a currency
// @Param id path string true "currencyID"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Currency}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /currency/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	currencyID := ctx.Param("id")
	currency, err := h.service.Delete(currencyID)
	if err != nil {
		shared.LogError("error deleting currency", LogHandler, "Delete", err, currencyID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCurrencyDeleting))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(currency))
}
