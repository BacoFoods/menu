package store

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler = "pkg/store/handler"

type RequestStoreCreate struct {
	Name      string  `json:"name"`
	BrandID   *uint   `json:"brand_id"`
	Enabled   bool    `json:"enabled"`
	Image     string  `json:"image,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// Handler to handle requests to the store service
type Handler struct {
	service Service
}

// NewHandler to create a new handler for the store service
func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Find to handle a request to find all stores
// @Tags Store
// @Summary To find stores
// @Description To find stores
// @Param name query string false "store name"
// @Param brand-id query string false "brand id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Store}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /store [get]
func (h Handler) Find(c *gin.Context) {
	query := make(map[string]string)

	name := c.Query("name")
	if name != "" {
		query["name"] = name
	}

	brandID := c.Query("brand-id")
	if brandID != "" {
		query["brand_id"] = brandID
	}

	stores, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding stores", LogHandler, "Find", err, stores)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingStore))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(stores))
}

// Get to handle a request to get a store
// @Tags Store
// @Summary To get a store
// @Description To get a store
// @Param id path string true "store id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Store}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /store/{id} [get]
func (h Handler) Get(c *gin.Context) {
	storeID := c.Param("id")

	store, err := h.service.Get(storeID)
	if err != nil {
		shared.LogError("error getting store", LogHandler, "Get", err, storeID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingStore))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(store))
}

// Create to handle a request to create a store
// @Tags Store
// @Summary To create a store
// @Description To create a store
// @Accept json
// @Produce json
// @Param store body object{name=string,brand_id=integer,enabled=boolean,image=string,latitude=number,longitude=number,address=string} true "store"
// @Success 200 {object} object{status=string,data=Store}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /store [post]
func (h Handler) Create(c *gin.Context) {
	var request Store
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogError("error binding json", LogHandler, "Create", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	store, err := h.service.Create(&request)
	if err != nil {
		shared.LogError("error creating store", LogHandler, "Create", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingStore))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(store))
}

// Update to handle a request to update a store
// @Tags Store
// @Summary To update a store
// @Description To update a store
// @Accept json
// @Produce json
// @Param id path string true "store id"
// @Param store body Store true "store"
// @Success 200 {object} object{status=string,data=Store}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /store/{id} [patch]
func (h Handler) Update(c *gin.Context) {
	var request Store
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Update", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	store, err := h.service.Update(&request)
	if err != nil {
		shared.LogError("error updating store", LogHandler, "Update", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingStore))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(store))
}

// Delete to handle a request to delete a store
// @Tags Store
// @Summary To delete a store
// @Description To delete a store
// @Accept json
// @Produce json
// @Param id path string true "store id"
// @Success 200 {object} object{status=string,data=Store}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /store/{id} [delete]
func (h Handler) Delete(c *gin.Context) {
	storeID := c.Param("id")
	store, err := h.service.Delete(storeID)
	if err != nil {
		shared.LogError("error deleting store", LogHandler, "Delete", err, storeID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingStore))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(store))
}

// AddChannel to handle a request to add a channel to a store
// @Tags Store
// @Summary To add a channel to a store
// @Description To add a channel to a store
// @Accept json
// @Produce json
// @Param id path string true "store id"
// @Param channelID path string true "channel id"
// @Success 200 {object} object{status=string,data=Store}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /store/{id}/channel/{channelID} [patch]
func (h Handler) AddChannel(c *gin.Context) {
	storeID := c.Param("id")
	channelID := c.Param("channelID")

	store, err := h.service.AddChannel(storeID, channelID)
	if err != nil {
		shared.LogError("error adding channel to store", LogHandler, "AddChannel", err, storeID, channelID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorAddingChannel))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(store))
}

// FindZonesByStore to handle a request to get zones by store
// @Tags Store
// @Summary To get zones by store
// @Description To get zones by store
// @Accept json
// @Produce json
// @Param id path string true "store id"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /store/{id}/zones [get]
func (h Handler) FindZonesByStore(c *gin.Context) {
	storeID := c.Param("id")

	zones, err := h.service.FindZonesByStore(storeID)
	if err != nil {
		shared.LogError("error getting zones by store", LogHandler, "GetZonesByStore", err, storeID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZonesGettingByStoreID))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(zones))
}

// GetZoneByStore to handle a request to get tables by store
// @Tags Store
// @Summary To get tables by store
// @Description To get tables by store
// @Accept json
// @Produce json
// @Param id path string true "store id"
// @Param zoneID path string true "zone id"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /store/{id}/zones/{zoneID} [get]
func (h Handler) GetZoneByStore(c *gin.Context) {
	storeID := c.Param("id")
	zoneID := c.Param("zoneID")

	zone, err := h.service.GetZoneByStore(storeID, zoneID)
	if err != nil {
		shared.LogError("error getting zone by store", LogHandler, "GetZoneByStore", err, storeID, zoneID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZonesGettingByStoreID))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(zone))
}
