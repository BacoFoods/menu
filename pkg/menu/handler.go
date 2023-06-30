package menu

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/menu/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Find to handle a request to find all menus
// @Tags Menu
// @Summary To find menus
// @Description To find menus
// @Param name query string false "menu name"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /menu [get]
func (h *Handler) Find(c *gin.Context) {
	query := make(map[string]string)
	name := c.Query("name")
	if name != "" {
		query["name"] = c.Query("name")
	}
	menus, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding menus", LogHandler, "Find", err, menus)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingMenu))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(menus))
}

// Get to handle a request to get a menu
// @Tags Menu
// @Summary To get a menu
// @Description To get a menu
// @Param id path string true "menu id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /menu/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	menuID := c.Param("id")
	menu, err := h.service.Get(menuID)
	if err != nil {
		shared.LogError("error getting menu", LogHandler, "Get", err, menu)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingMenu))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(menu))
}

// Create to handle a request to create a menu
// @Tags Menu
// @Summary To create a menu
// @Description To create a menu
// @Accept json
// @Produce json
// @Param menu body Menu true "menu"
// @Success 200 {object} object{status=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /menu [post]
func (h *Handler) Create(c *gin.Context) {
	var requestBody Menu
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request body", LogHandler, "Create", err, requestBody)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	menu, err := h.service.Create(&requestBody)
	if err != nil {
		shared.LogError("error creating menu", LogHandler, "Create", err, requestBody)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingMenu))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(menu))
}

// Update to handle a request to update a menu
// @Tags Menu
// @Summary To update a menu
// @Description To update a menu
// @Accept json
// @Produce json
// @Param id path string true "menu id"
// @Param menu body Menu true "menu"
// @Success 200 {object} object{status=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /menu/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	var requestBody Menu
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request body", LogHandler, "Update", err, requestBody)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	menu, err := h.service.Update(&requestBody)
	if err != nil {
		shared.LogError("error updating menu", LogHandler, "Update", err, requestBody)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingMenu))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(menu))
}

// Delete to handle a request to delete a menu
// @Tags Menu
// @Summary To delete a menu
// @Description To delete a menu
// @Accept json
// @Produce json
// @Param id path string true "menu id"
// @Success 200 {object} object{status=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /menu/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	menuID := c.Param("id")
	menu, err := h.service.Delete(menuID)
	if err != nil {
		shared.LogError("error deleting menu", LogHandler, "Delete", err, menu)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingMenu))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(menu))
}
