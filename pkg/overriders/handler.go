package overriders

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/overriders/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Find to handle a request to find all overriders
// @Tags Overriders
// @Summary To find overriders
// @Description To find overriders
// @Param name query string false "overriders name"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Overriders}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /overriders [get]
func (h *Handler) Find(c *gin.Context) {
	query := make(map[string]string)
	name := c.Query("name")
	if name != "" {
		query["name"] = c.Query("name")
	}
	overriders, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding overriders", LogHandler, "Find", err, overriders)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingOverriders))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(overriders))
}

// Get to handle a request to get an overriders
// @Tags Overriders
// @Summary To get an overriders
// @Description To get an overriders
// @Param id path string true "overriders id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Overriders}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /overriders/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	overridersID := c.Param("id")
	overriders, err := h.service.Get(overridersID)
	if err != nil {
		shared.LogError("error getting overriders", LogHandler, "Get", err, overriders)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingOverriders))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(overriders))
}

// Create to handle a request to create an overriders
// @Tags Overriders
// @Summary To create an overriders
// @Description To create an overriders
// @Accept json
// @Produce json
// @Param overriders body Overriders true "overriders"
// @Success 200 {object} object{status=string,data=Overriders}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /overriders [post]
func (h *Handler) Create(c *gin.Context) {
	var request Overriders
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding overriders", LogHandler, "Create", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	overriders, err := h.service.Create(&request)
	if err != nil {
		shared.LogError("error creating overriders", LogHandler, "Create", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingOverriders))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(overriders))
}

// Update to handle a request to update an overriders
// @Tags Overriders
// @Summary To update an overriders
// @Description To update an overriders
// @Accept json
// @Produce json
// @Param id path string true "overriders id"
// @Param overriders body Overriders true "overriders"
// @Success 200 {object} object{status=string,data=Overriders}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /overriders/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	var request Overriders
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding overriders", LogHandler, "Update", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	overriders, err := h.service.Update(&request)
	if err != nil {
		shared.LogError("error updating overriders", LogHandler, "Update", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingOverriders))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(overriders))
}

// Delete to handle a request to delete an overriders
// @Tags Overriders
// @Summary To delete an overriders
// @Description To delete an overriders
// @Param id path string true "overriders id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Overriders}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /overriders/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	overridersID := c.Param("id")
	overriders, err := h.service.Delete(overridersID)
	if err != nil {
		shared.LogError("error deleting overriders", LogHandler, "Delete", err, overriders)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingOverriders))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(overriders))
}
