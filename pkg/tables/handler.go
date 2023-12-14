package tables

import (
	"net/http"
	"strconv"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const (
	LogHandler = "pkg/tables/handler"
)

type RequestZoneCreate struct {
	Name        string `json:"name" binding:"required"`
	StoreID     uint   `json:"store_id" binding:"required"`
	Active      bool   `json:"active"`
	TableNumber int    `json:"table_number"`
	TableAmount int    `json:"table_amount"`
}

type RequestZoneUpdate struct {
	Name    string `json:"name"`
	StoreID uint   `json:"store_id"`
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

// Get to handle get tables request
// @Tags Tables
// @Summary Get tables
// @Description Get tables
// @Param id path string true "table id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Table}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /tables/{id} [get]
func (h Handler) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	table, err := h.service.Get(id)
	if err != nil {
		shared.LogError("error getting table", LogHandler, "Get", err, table)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTableGetting))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(table))
}

// Release table
// @Tags Tables
// @Summary Release tables
// @Description Release tables
// @Param id path string true "table id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Table}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /tables/{id}/release [post]
func (h Handler) Release(ctx *gin.Context) {
	pid := ctx.Param("id")
	id, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		shared.LogError("error parsing table id", LogHandler, "Release", err, pid)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTableGetting))
		return
	}

	table, err := h.service.ReleaseTable(uint(id))
	if err != nil {
		shared.LogError("error releasing table", LogHandler, "Release", err, id)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTableUpdating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(table))
}

// Find to handle find tables request
// @Tags Tables
// @Summary Find tables
// @Description Find tables
// @Param name query string false "table name"
// @Param zoneID query int false "zone id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Table}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /tables [get]
func (h Handler) Find(ctx *gin.Context) {
	query := make(map[string]any)

	name := ctx.Query("name")
	if name != "" {
		query["name"] = name
	}

	zoneID := ctx.Query("zone-id")
	if zoneID != "" {
		query["zone_id"] = zoneID
	}

	tables, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding tables", LogHandler, "Find", err, tables)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTableFinding))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(tables))
}

// Create to handle create table request
// @Tags Tables
// @Summary Create table
// @Description Create table
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param table body Table true "table"
// @Success 200 {object} object{status=string,data=Table}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /tables [post]
func (h Handler) Create(ctx *gin.Context) {
	var table Table
	if err := ctx.ShouldBindJSON(&table); err != nil {
		shared.LogError("error binding table", LogHandler, "Create", err, table)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	newTable, err := h.service.Create(&table)
	if err != nil {
		shared.LogError("error creating table", LogHandler, "Create", err, table)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTableCreating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(newTable))
}

// Update to handle update table request
// @Tags Tables
// @Summary Update table
// @Description Update table
// @Param id path string true "table id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Table}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /tables/{id} [patch]
func (h Handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var table Table
	if err := ctx.ShouldBindJSON(&table); err != nil {
		shared.LogError("error binding table", LogHandler, "Update", err, table)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	updateTable, err := h.service.Update(id, &table)
	if err != nil {
		shared.LogError("error updating table", LogHandler, "Update", err, table)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTableUpdating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(updateTable))
}

// Delete to handle delete table request
// @Tags Tables
// @Summary Delete table
// @Description Delete table
// @Param id path string true "table id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /tables/{id} [delete]
func (h Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.service.Delete(id); err != nil {
		shared.LogError("error deleting table", LogHandler, "Delete", err, id)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTableDeleting))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(nil))
}

// ScanQR to handle scan qr request
// @Tags Tables
// @Summary Scan QR
// @Description Scan QR
// @Param qrId path string true "qr id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Table}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /public/tables/scan/{qrId} [get]
func (h Handler) ScanQR(ctx *gin.Context) {
	qrID := ctx.Param("qrId")

	table, err := h.service.ScanQR(qrID)
	if err != nil {
		shared.LogError("error scanning qr", LogHandler, "ScanQR", err, qrID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTableScanningQR))
		return
	}

	if table == nil {
		ctx.JSON(http.StatusNotFound, shared.ErrorResponse(ErrorTableNotFound))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(table))
}

// GenerateQR to handle qr generation request
// @Tags Tables
// @Summary Generate QR
// @Description Generate QR
// @Param id path string true "table id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Table}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /tables/{id}/generate [post]
func (h Handler) GenerateQR(ctx *gin.Context) {
	id := ctx.Param("id")

	table, err := h.service.GenerateQR(id)
	if err != nil {
		shared.LogError("error generating qr", LogHandler, "GenerateQR", err, id)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorTableGeneratingQR))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(table))
}

// Find to handle a request to find all zones
// @Tags Zones
// @Summary To find zones
// @Description To find zones
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "Name"
// @Param storeID query int false "Store ID"
// @Success 200 {object} object{status=string,data=[]Zone}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /zone [get]
func (h Handler) FindZones(c *gin.Context) {
	filters := make(map[string]any)
	zoneName := c.Query("name")
	if zoneName != "" {
		filters["name"] = zoneName
	}

	storeID := c.Query("storeID")
	if storeID != "" {
		filters["store_id"] = storeID
	}

	zones, err := h.service.FindZones(filters)
	if err != nil {
		shared.LogError("error finding zones", LogHandler, "FindZones", err, zones)
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
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Success 200 {object} object{status=string,data=Zone}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /zone/{id} [get]
func (h Handler) GetZone(c *gin.Context) {
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
// @Security ApiKeyAuth
// @Param request body RequestZoneCreate true "Zone"
// @Success 200 {object} object{status=string,data=Zone}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /zone [post]
func (h Handler) CreateZone(c *gin.Context) {
	var req RequestZoneCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Create", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorZoneBadRequest))
		return
	}

	zone := Zone{
		Active:  req.Active,
		Name:    req.Name,
		StoreID: &req.StoreID,
	}

	newZone, err := h.service.CreateZone(&zone, req.TableNumber, req.TableAmount)
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
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Param request body RequestZoneUpdate true "Zone"
// @Success 200 {object} object{status=string,data=Zone}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /zone/{id} [patch]
func (h Handler) UpdateZone(c *gin.Context) {
	id := c.Param("id")
	var req RequestZoneUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Update", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorZoneBadRequest))
		return
	}

	zone := Zone{
		Name:    req.Name,
		StoreID: &req.StoreID,
	}

	updatedZone, err := h.service.UpdateZone(id, &zone)
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
// @Security ApiKeyAuth
// @Param id path int true "ID"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /zone/{id} [delete]
func (h Handler) DeleteZone(c *gin.Context) {
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
// @Security ApiKeyAuth
// @Param id path int true "Zone ID"
// @Param request body RequestAddTable true "Add Table"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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
// @Security ApiKeyAuth
// @Param id path int true "Zone ID"
// @Param request body RequestAddTable true "Remove Table"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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

// Enable to handle a request to enable/disable a zone
// @Tags Zones
// @Summary To enable/disable a zone
// @Description To enable/disable a zone
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Zone ID"
// @Success 200 {object} object{status=string,data=Zone}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /zone/{id}/enable [patch]
func (h Handler) Enable(c *gin.Context) {
	id := c.Param("id")

	zone, err := h.service.EnableZone(id)
	if err != nil {
		shared.LogError("error enabling zone", LogHandler, "Enable", err, id)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorZoneEnabling))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(zone))
}
