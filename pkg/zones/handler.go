package zones

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	LogHandler = "pkg/zone/handler"
)

type RequestZoneCreate struct {
	Name        string `json:"name" binding:"required"`
	StoreID     uint   `json:"store_id" binding:"required"`
	TableNumber int    `json:"table_number"`
	TableAmount int    `json:"table_amount"`
}

type RequestAddTable struct {
	Tables []uint `json:"tables" binding:"required"`
}

type RequestRemoveTable struct {
	Tables []uint `json:"tables" binding:"required"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Find to handle a request to find all zones
// @Tags Zones
// @Summary To find zones
// @Description To find zones
// @Accept json
// @Produce json
// @Param name query string false "Name"
// @Param storeID query int false "Store ID"
// @Success 200 {object} object{status=string,data=[]Zone}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /zone [get]
func (h Handler) Find(c *gin.Context) {
	filters := make(map[string]any)
	zoneName := c.Query("name")
	if zoneName != "" {
		filters["name"] = zoneName
	}

	storeID := c.Query("storeID")
	if storeID != "" {
		filters["store_id"] = storeID
	}

	zones, err := h.service.Find(filters)
	if err != nil {
		shared.LogError("error finding zones", LogHandler, "Find", err, zones)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZoneFinding))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(zones))
}

// Get to handle a request to get a zone
// @Tags Zones
// @Summary To get a zone
// @Description To get a zone
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} object{status=string,data=Zone}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /zone/{id} [get]
func (h Handler) Get(c *gin.Context) {
	id := c.Param("id")
	zone, err := h.service.Get(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZoneFinding))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(zone))
}

// Create to handle a request to create a zone
// @Tags Zones
// @Summary To create a zone
// @Description To create a zone
// @Accept json
// @Produce json
// @Param request body RequestZoneCreate true "Zone"
// @Success 200 {object} object{status=string,data=Zone}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /zone [post]
func (h Handler) Create(c *gin.Context) {
	var req RequestZoneCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Create", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorZoneBadRequest))
		return
	}

	zone := Zone{
		Name:    req.Name,
		StoreID: &req.StoreID,
	}

	newZone, err := h.service.Create(&zone, req.TableNumber, req.TableAmount)
	if err != nil {
		shared.LogError("error creating zone", LogHandler, "Create", err, newZone)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZoneCreating))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(newZone))
}

// Update to handle a request to update a zone
// @Tags Zones
// @Summary To update a zone
// @Description To update a zone
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param request body RequestZoneCreate true "Zone"
// @Success 200 {object} object{status=string,data=Zone}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /zone/{id} [patch]
func (h Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req RequestZoneCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Update", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorZoneBadRequest))
		return
	}

	zone := Zone{
		Name:    req.Name,
		StoreID: &req.StoreID,
	}

	updatedZone, err := h.service.Update(id, &zone)
	if err != nil {
		shared.LogError("error updating zone", LogHandler, "Update", err, updatedZone)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZoneUpdating))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(updatedZone))
}

// Delete to handle a request to delete a zone
// @Tags Zones
// @Summary To delete a zone
// @Description To delete a zone
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /zone/{id} [delete]
func (h Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.service.Delete(id)
	if err != nil {
		shared.LogError("error deleting zone", LogHandler, "Delete", err, id)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZoneDeleting))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(nil))
}

// AddTables to handle a request to add tables to a zone
// @Tags Zones
// @Summary To add tables to a zone
// @Description To add tables to a zone
// @Accept json
// @Produce json
// @Param id path int true "Zone ID"
// @Param request body RequestAddTable true "Add Table"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /zone/{id}/tables/add [patch]
func (h Handler) AddTables(c *gin.Context) {
	id := c.Param("id")
	var req RequestAddTable
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "AddTablesToZone", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorZoneBadRequest))
		return
	}

	err := h.service.AddTables(id, req.Tables)
	if err != nil {
		shared.LogError("error adding tables to zone", LogHandler, "AddTablesToZone", err, id, req.Tables)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZoneAddingTables))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse("tables added to zone"))
}

// RemoveTables to handle a request to remove tables from a zone
// @Tags Zones
// @Summary To remove tables from a zone
// @Description To remove tables from a zone
// @Accept json
// @Produce json
// @Param id path int true "Zone ID"
// @Param request body RequestAddTable true "Remove Table"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /zone/{id}/tables/remove [patch]
func (h Handler) RemoveTables(c *gin.Context) {
	id := c.Param("id")

	var req RequestRemoveTable
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "RemoveTablesFromZone", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorZoneBadRequest))
		return
	}

	err := h.service.RemoveTables(id, req.Tables)
	if err != nil {
		shared.LogError("error removing tables from zone", LogHandler, "RemoveTablesFromZone", err, id, req.Tables)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZoneRemovingTables))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse("tables removed from zone"))
}
