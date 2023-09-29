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
// @Param storeID query string false "Store ID"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Surcharge}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /surcharge [get]
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

	storeID := c.Query("storeID")
	if storeID != "" {
		query["store_id"] = storeID
	}

	surcharges, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding surcharges", LogHandler, "Find", err, query)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(surcharges))
}

// Get surcharge
// @Tags Surcharge
// @Summary Get surcharge
// @Description Get surcharge
// @Param id path string true "Surcharge ID"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Surcharge}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /surcharge/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")

	surcharge, err := h.service.Get(id)
	if err != nil {
		shared.LogError("error getting surcharge", LogHandler, "Get", err, id)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(surcharge))
}

// Create surcharge
// @Tags Surcharge
// @Summary Create surcharge
// @Description Create surcharge
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param surcharge body Surcharge true "Surcharge"
// @Success 200 {object} object{status=string,data=Surcharge}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /surcharge [post]
func (h *Handler) Create(c *gin.Context) {
	var surcharge Surcharge
	if err := c.ShouldBindJSON(&surcharge); err != nil {
		shared.LogError("error binding surcharge", LogHandler, "Create", err, surcharge)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	newSurcharge, err := h.service.Create(&surcharge)
	if err != nil {
		shared.LogError("error creating surcharge", LogHandler, "Create", err, surcharge)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(newSurcharge))
}

// Update surcharge
// @Tags Surcharge
// @Summary Update surcharge
// @Description Update surcharge
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Surcharge ID"
// @Param surcharge body Surcharge true "Surcharge"
// @Success 200 {object} object{status=string,data=Surcharge}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /surcharge/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var surcharge Surcharge
	if err := c.ShouldBindJSON(&surcharge); err != nil {
		shared.LogError("error binding surcharge", LogHandler, "Update", err, surcharge)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	newSurcharge, err := h.service.Update(id, &surcharge)
	if err != nil {
		shared.LogError("error updating surcharge", LogHandler, "Update", err, surcharge)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(newSurcharge))
}

// Delete surcharge
// @Tags Surcharge
// @Summary Delete surcharge
// @Description Delete surcharge
// @Param id path string true "Surcharge ID"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Surcharge}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /surcharge/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	surcharge, err := h.service.Delete(id)
	if err != nil {
		shared.LogError("error deleting surcharge", LogHandler, "Delete", err, id)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(surcharge))
}
