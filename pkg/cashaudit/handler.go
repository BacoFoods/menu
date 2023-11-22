package cashaudit

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	LogHandler = "cashaudit/handler.go"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// OrdersClosedValidation to handle orders closed validation request
// @Tags Cash Audit
// @Summary To validate if all orders are closed
// @Description To validate if all orders are closed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=string}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /cash-audit/orders-closed [get]
func (h Handler) OrdersClosedValidation(c *gin.Context) {
	storeID, ok := c.Get("store_id")
	if !ok {
		shared.LogWarn("error getting store id", LogHandler, "OrdersClosedValidation", fmt.Errorf(ErrorCashAuditGettingStoreID))
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCashAuditGettingStoreID))
	}

	if err := h.service.AllOrdersClosed(storeID.(string)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse("all orders closed"))
}

// Get to handle get cash audit request
// @Tags Cash Audit
// @Summary Just calculate the cash audit only for closed orders
// @Description Just calculate the cash audit only for closed orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=DTOCashAuditCategories}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /cash-audit [get]
func (h Handler) Get(c *gin.Context) {
	storeID, ok := c.Get("store_id")
	if !ok {
		shared.LogWarn("error getting store id", LogHandler, "Get", fmt.Errorf(ErrorCashAuditGettingStoreID))
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCashAuditGettingStoreID))
	}

	cashAudit, err := h.service.Get(storeID.(string))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(ToDTOCashAuditCategories(*cashAudit)))
}

// Create to handle create cash audit request
// @Tags Cash Audit
// @Summary To create cash audit only for closed orders
// @Description To create cash audit only for closed orders, if cash audit already exists, it will return the existing cash audit
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param cashAudit body DTOCashAudit true "Cash Audit"
// @Success 200 {object} object{status=string,data=DTOCashAuditCategories}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /cash-audit [post]
func (h Handler) Create(c *gin.Context) {
	storeID, ok := c.Get("store_id")
	if !ok {
		shared.LogWarn("error getting store id", LogHandler, "Create", fmt.Errorf(ErrorCashAuditGettingStoreID))
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCashAuditGettingStoreID))
	}

	accountID, ok := c.Get("account_id")
	if !ok {
		shared.LogWarn("error getting account id", LogHandler, "Create", fmt.Errorf(ErrorCashAuditGettingStoreID))
	}
	u64, err := strconv.ParseUint(accountID.(string), 10, 32)
	if err != nil {
		shared.LogWarn("error parsing account id", LogHandler, "Create", err)
	}
	accountUintID := uint(u64)

	var dtoCashAudit DTOCashAudit
	if err := c.ShouldBindJSON(&dtoCashAudit); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	cashAudit := dtoCashAudit.ToCashAudit()
	cashAudit.CashierAccountID = &accountUintID
	createdCashAudit, err := h.service.Create(storeID.(string), &cashAudit)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(ToDTOCashAuditCategories(*createdCashAudit)))
}

// Confirm to handle confirm cash audit request
// @Tags Cash Audit
// @Summary To confirm cash audit only for closed orders
// @Description To confirm cash audit only for closed orders, if cash audit already exists, it will return the existing cash audit
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param cashAudit body DTOCashAuditConfirmationRequest true "Cash Audit Confirmation"
// @Success 200 {object} object{status=string,data=CashAudit}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /cash-audit/confirm [post]
func (h Handler) Confirm(c *gin.Context) {
	var dto DTOCashAuditConfirmationRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	cashAudit, err := h.service.Confirm(dto.CashAuditID, dto.Observations)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(&cashAudit))
}
