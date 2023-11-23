package facturacion

import (
	"net/http"
	"strconv"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const LogHandler string = "pkg/facturacion/handler"

type Handler struct {
	service *FacturacionService
}

// NewHandler to create a new handler
func NewHandler(service *FacturacionService) *Handler {
	return &Handler{service}
}

// CreateConfig to handle a request for create a facturacion config
// @Tags Facturacion
// @Summary To create a facturacion config
// @Description To create a facturacion config
// @Param storeId path string true "store id"
// @Param config body FacturacionConfig true "facturacion config request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=FacturacionConfig}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /store/{storeId}/facturacion/config [post]
func (h *Handler) CreateConfig(c *gin.Context) {
	id := c.Param("id")
	storeId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "CreateConfig", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse("invalid store id"))
		return
	}

	var requestBody FacturacionConfig
	if err := c.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "CreateConfig", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse("invalid request body"))
		return
	}

	requestBody.StoreID = uint(storeId)

	config, err := h.service.CreateConfig(&requestBody)
	if err != nil {
		shared.LogError("error creating facturacion config", LogHandler, "CreateConfig", err, config)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse("error creating facturacion config"))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(config))
}

// FindConfig to handle a request for find a facturacion config
// @Tags Facturacion
// @Summary To find a facturacion config
// @Description To find a facturacion config
// @Param id path string true "store id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=FacturacionConfig}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /store/{id}/facturacion/config [get]
func (h *Handler) FindConfig(c *gin.Context) {
	id := c.Param("id")
	storeId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "FindConfig", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse("invalid store id"))
		return
	}

	config, err := h.service.FindConfig(uint(storeId))
	if err != nil {
		shared.LogError("error finding facturacion config", LogHandler, "FindConfig", err, config)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse("error finding facturacion config"))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(config))
}

// UpdateConfig to handle a request for update a facturacion config
// @Tags Facturacion
// @Summary To update a facturacion config
// @Description To update a facturacion config
// @Param storeId path string true "store id"
// @Param id path string true "config id"
// @Param config body FacturacionConfig true "facturacion config request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=FacturacionConfig}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /store/{storeId}/facturacion/config/{id} [put]
func (h *Handler) UpdateConfig(c *gin.Context) {
	id := c.Param("id")
	storeId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "UpdateConfig", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse("invalid store id"))
		return
	}

	id = c.Param("configId")
	configId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "UpdateConfig", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse("invalid config id"))
		return
	}

	var requestBody FacturacionConfig
	if err := c.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "UpdateConfig", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse("invalid request body"))
		return
	}

	requestBody.ID = uint(configId)
	requestBody.StoreID = uint(storeId)

	config, err := h.service.UpdateConfig(&requestBody)
	if err != nil {
		shared.LogError("error updating facturacion config", LogHandler, "UpdateConfig", err, config)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse("error updating facturacion config"))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(config))
}
