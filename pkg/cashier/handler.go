package cashier

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

// Open to handle open cashier request
// @Tags Cashier
// @Summary Open cashier
// @Description Open cashier
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /cashier/open [post]
func (h *Handler) Open(c *gin.Context) {
	if err := h.service.Open(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse("cashier opened successfully"))
}

// Close to handle close cashier request
// @Tags Cashier
// @Summary Close cashier
// @Description Close cashier
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string}
// @Failure 401 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /cashier/close [post]
func (h *Handler) Close(c *gin.Context) {
	if err := h.service.Close(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse("cashier closed successfully"))
}
