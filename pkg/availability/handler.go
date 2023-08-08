package availability

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const LogHandler string = "pkg/availability/handler"

type ResponsePlaces struct {
	Entity   string `json:"entity"`
	EntityID uint   `json:"entity_id"`
	Place    string `json:"place"`
	Places   []any  `json:"places"`
}

type EnableRequest struct {
	Enable bool `json:"enable"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// RemoveEntity to handle a request to remove a menu or category
// @Summary Remove menu or category
// @Description Remove menu or category
// @Tags Availability
// @Produce json
// @Param entity path string true "Entity"
// @Param entity-id path string true "Entity ID"
// @Param place path string true "Place"
// @Param place-id path string true "Place ID"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /availability/{entity}/{entity-id}/{place}/{place-id} [delete]
func (h *Handler) RemoveEntity(c *gin.Context) {
	entity, err := GetEntity(c.Param("entity"))
	if err != nil {
		shared.LogWarn("warning getting entity", LogHandler, "RemoveEntity", err, entity)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	place, err := GetPlace(c.Param("place"))
	if err != nil {
		shared.LogWarn("warning getting place", LogHandler, "RemoveEntity", err, entity, place)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	entityID, err := strconv.ParseUint(c.Param("entity-id"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing entity id", LogHandler, "RemoveEntity", err, entity, entityID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	placeID, err := strconv.ParseUint(c.Param("place-id"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing place id", LogHandler, "RemoveEntity", err, entity, entityID, place, placeID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	if err := h.service.RemoveEntity(entity, place, uint(entityID), uint(placeID)); err != nil {
		shared.LogError("error removing entity", LogHandler, "RemoveEntity", err, entity, entityID, place, placeID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorRemovingEntity))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(fmt.Sprintf("%v: %v in %v: %v", entity, entityID, place, placeID)))
}

// EnableEntity to handle a request to enable a menu or category
// @Summary Enable menu or category
// @Description Enable menu or category
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
// @Router /availability/{entity}/{entity-id}/{place}/{place-id} [put]
func (h *Handler) EnableEntity(c *gin.Context) {
	entity, err := GetEntity(c.Param("entity"))
	if err != nil {
		shared.LogWarn("warning getting entity", LogHandler, "EnableEntity", err, entity)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	place, err := GetPlace(c.Param("place"))
	if err != nil {
		shared.LogWarn("warning getting place", LogHandler, "EnableEntity", err, entity, place)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	entityID, err := strconv.ParseUint(c.Param("entity-id"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing entity id", LogHandler, "EnableEntity", err, entity, entityID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	placeID, err := strconv.ParseUint(c.Param("place-id"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing place id", LogHandler, "EnableEntity", err, entity, entityID, place, placeID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	var body EnableRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogWarn("warning binding body", LogHandler, "EnableEntity", err, entity, entityID, place, placeID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	if err := h.service.EnableEntity(entity, place, uint(entityID), uint(placeID), body.Enable); err != nil {
		shared.LogError("error enabling entity", LogHandler, "EnableEntity", err, entity, entityID, place, placeID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorEnablingEntity))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(fmt.Sprintf("%v: %v in %v: %v enable: %v", entity, entityID, place, placeID, body.Enable)))
}

// FindEntities to handle a request to find all entities
// @Summary Find all entities
// @Description Find all entities
// @Tags Availability
// @Produce json
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /availability/entities [get]
func (h *Handler) FindEntities(c *gin.Context) {
	entities := h.service.FindEntities()
	c.JSON(http.StatusOK, shared.SuccessResponse(entities))
}

// FindPlaces to handle a request to find all places
// @Summary Find all places
// @Description Find all places
// @Tags Availability
// @Produce json
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /availability/places [get]
func (h *Handler) FindPlaces(c *gin.Context) {
	places := h.service.FindPlaces()
	c.JSON(http.StatusOK, shared.SuccessResponse(places))
}

// Find to handle a request to find availability for entity and place
// @Summary Find availability for entity and place
// @Description Find availability for entity and place
// @Tags Availability
// @Accept json
// @Produce json
// @Param entity path string true "Entity"
// @Param entity-id path string true "Entity ID"
// @Param place path string true "Place"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /availability/{entity}/{entity-id}/{place} [get]
func (h *Handler) Find(c *gin.Context) {
	entity, err := GetEntity(c.Param("entity"))
	if err != nil {
		shared.LogWarn("warning getting entity", LogHandler, "Find", err, entity)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	place, err := GetPlace(c.Param("place"))
	if err != nil {
		shared.LogWarn("warning getting place", LogHandler, "Find", err, entity, place)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	entityID, err := strconv.ParseUint(c.Param("entity-id"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing entity id", LogHandler, "Find", err, entity, entityID)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	places, err := h.service.Find(entity, place, uint(entityID))
	if err != nil {
		shared.LogError("error finding availability", LogHandler, "Find", err, entity, entityID, place)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingAvailability))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(ResponsePlaces{
		Entity:   string(entity),
		EntityID: uint(entityID),
		Place:    string(place),
		Places:   places,
	}))
}

// Get to handle a request to get availability for entity and place
// @Summary Get availability for entity and place
// @Description Get availability for entity and place
// @Tags Availability
// @Accept json
// @Produce json
// @Param entity path string true "Entity"
// @Param entity-id path string true "Entity ID"
// @Param place path string true "Place"
// @Param place-id path string true "Place ID"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /availability/{entity}/{entity-id}/{place}/{place-id} [get]
func (h *Handler) Get(c *gin.Context) {
	entity, err := GetEntity(c.Param("entity"))
	if err != nil {
		shared.LogWarn("warning getting entity", LogHandler, "Get", err, c.Param("entity"))
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	place, err := GetPlace(c.Param("place"))
	if err != nil {
		shared.LogWarn("warning getting place", LogHandler, "Get", err, c.Param("place"))
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	entityID, err := strconv.ParseUint(c.Param("entity-id"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing entity id", LogHandler, "Get", err, c.Param("entity-id"))
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	placeID, err := strconv.ParseUint(c.Param("place-id"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing place id", LogHandler, "Get", err, c.Param("place-id"))
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	placeItem, err := h.service.Get(entity, place, uint(entityID), uint(placeID))
	if err != nil {
		shared.LogError("error getting availability", LogHandler, "Get", err, entity, entityID, place, placeID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingAvailability))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(ResponsePlaces{
		Entity:   string(entity),
		EntityID: uint(entityID),
		Place:    string(place),
		Places:   []any{placeItem},
	}))
}
