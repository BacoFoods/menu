package tables

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	LogHandler = "pkg/tables/handler"
)

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
// @Failure 403 {object} shared.Response
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
// @Failure 403 {object} shared.Response
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
// @Failure 403 {object} shared.Response
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
// @Failure 403 {object} shared.Response
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
// @Failure 403 {object} shared.Response
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
