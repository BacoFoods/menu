package store

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler = "pkg/store/handler"

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
		query["name"] = c.Query("name")
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
// @Param store body Store true "store"
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
