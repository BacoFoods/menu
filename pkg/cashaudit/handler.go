package cashaudit

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
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

// Get to handle get cash audit request
// @Tags Cash Audit
// @Summary To get cash audit
// @Description To get cash audit
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=CashAudit}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /cash-audit [get]
func (h Handler) Get(c *gin.Context) {
	storeID, ok := c.Get("store_id")
	if !ok {
		shared.LogWarn("error getting store id", LogHandler, "Get", fmt.Errorf(ErrorGettingStoreID))
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorGettingStoreID))
	}

	cashAudit, err := h.service.Get(storeID.(string))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingCashAudit))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(&cashAudit))
}
