package surcharge

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler = "pkg/surcharge/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Find surcharges
// @Tags Surcharge
// @Summary Find surcharges
// @Description Find surcharges
// @Param brandID query string false "Brand ID"
// @Param name query string false "Name"
// @Accept json
// @Produce json
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /surcharges [get]
func (h *Handler) Find(c *gin.Context) {
	query := make(map[string]string)

	name := c.Query("name")
	if name != "" {
		query["name"] = name
	}

	brandID := c.Query("brandID")
	if brandID != "" {
		query["brand_id"] = brandID
	}

	surcharges, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding surcharges", LogHandler, "Find", err, query)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(surcharges))
}
