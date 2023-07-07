package availability

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/availability/handler"

type EnableRequest struct {
	Enable bool `json:"enable"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// EnableEntity to handle a request to enable a menu or category
// @Summary Enable menu
// @Description Enable menu
// @Tags Availability
// @Produce json
// @Param entity path string true "Entity"
// @Param entity-id path string true "Entity ID"
// @Param place path string true "Place"
// @Param place-id path string true "Place ID"
// @Param enable body EnableRequest true "Enable"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /availability/{entity}/{entity-id}/{place}/{place-id} [post]
func (h *Handler) EnableEntity(c *gin.Context) {
	entity := c.Param("entity")
	entityID := c.Param("entity-id")
	place := c.Param("place")
	placeID := c.Param("place-id")

	var body EnableRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogWarn("warning binding body", LogHandler, "EnableEntity", err, entity, entityID, place, placeID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	err := h.service.EnableEntity(entity, entityID, place, placeID, body.Enable)
	if err != nil {
		shared.LogError("error enabling entity", LogHandler, "EnableEntity", err, entity, entityID, place, placeID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorEnablingEntity))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(fmt.Sprintf("%s: %s in %s: %s enable: %v", entity, entityID, place, placeID, body.Enable)))
}
