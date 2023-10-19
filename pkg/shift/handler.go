package shift

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// OpenShift to handle open cashier shift request
// @Tags Shift
// @Summary Open cashier shift
// @Description Open cashier shift
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /shift-cashier/open [post]
func (h *Handler) OpenShift(c *gin.Context) {
	if err := h.service.OpenShift(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse("cashier opened successfully"))
}

// CloseShift to handle close cashier shift request
// @Tags Shift
// @Summary Close cashier shift
// @Description Close cashier shift
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /shift-cashier/close [post]
func (h *Handler) CloseShift(c *gin.Context) {
	if err := h.service.CloseShift(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse("cashier closed successfully"))
}
