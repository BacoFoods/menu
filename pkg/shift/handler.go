package shift

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	LogHandler string = "pkg/shift/handler"
)

type Handler struct {
	service Service
}

type RequestOpenShift struct {
	StartBalance float64 `json:"start_balance" binding:"required"`
}

type RequestCloseShift struct {
	EndBalance float64 `json:"end_balance" binding:"required"`
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Open to handle open shift request
// @Tags Shift
// @Summary Open shift
// @Description Open shift
// @Accept json
// @Produce json
// @Param openShift body RequestOpenShift true "Request body"
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /shift/open [post]
func (h *Handler) Open(c *gin.Context) {
	accountID, ok := c.Get("account_uuid")
	if !ok {
		shared.LogError("account_id not found from jwt", LogHandler, "Open", nil)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse("account_id not found from jwt"))
		return
	}

	var req RequestOpenShift
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	shift, err := h.service.Open(accountID.(string), req.StartBalance)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(shift))
}

// Close to handle close shift request
// @Tags Shift
// @Summary Close shift
// @Description Close shift
// @Accept json
// @Produce json
// @Param closeShift body RequestCloseShift true "Request body"
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /shift/close [post]
func (h *Handler) Close(c *gin.Context) {
	accountID, ok := c.Get("account_uuid")
	if !ok {
		shared.LogError("account_id not found from jwt", LogHandler, "Open", nil)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse("account_id not found from jwt"))
		return
	}

	var req RequestCloseShift
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	shift, err := h.service.Close(accountID.(string), req.EndBalance)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(shift))
}
