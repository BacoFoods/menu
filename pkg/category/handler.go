package category

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/category/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Find to handle a request to find all categories
// @Tags Category
// @Summary To find categories
// @Description To find categories
// @Param name query string false "category name"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /category [get]
func (h *Handler) Find(c *gin.Context) {
	query := make(map[string]string)
	name := c.Query("name")
	if name != "" {
		query["name"] = c.Query("name")
	}
	categories, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding categories", LogHandler, "Find", err, categories)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingCategory))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(categories))
}

// Get to handle a request to get a category
// @Tags Category
// @Summary To get a category
// @Description To get a category
// @Param id path string true "category id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /category/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	categoryID := c.Param("id")
	category, err := h.service.Get(categoryID)
	if err != nil {
		shared.LogError("error getting category", LogHandler, "Get", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingCategory))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(category))
}

// Create to handle a request to create a category
// @Tags Category
// @Summary To create a category
// @Description To create a category
// @Param category body Category true "category request"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /category [post]
func (h *Handler) Create(c *gin.Context) {
	var requestBody Category
	if err := c.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	category, err := h.service.Create(&requestBody)
	if err != nil {
		shared.LogError("error creating category", LogHandler, "Create", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingCategory))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(category))
}

// Update to handle a request to update a category
// @Tags Category
// @Summary To update a category
// @Description To update a category
// @Param id path string true "category id"
// @Param category body Category true "category request"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /category/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	var requestBody Category
	if err := c.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Update", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}
	category, err := h.service.Update(&requestBody)
	if err != nil {
		shared.LogError("error updating category", LogHandler, "Update", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingCategory))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(category))
}

// Delete to handle a request to delete a category
// @Tags Category
// @Summary To delete a category
// @Description To delete a category
// @Param id path string true "category id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /category/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	categoryID := c.Param("id")
	category, err := h.service.Delete(categoryID)
	if err != nil {
		shared.LogError("error deleting category", LogHandler, "Delete", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingCategory))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(category))
}

// GetMenus to handle a request to get menus from a category
// @Tags Category
// @Summary To get menus from a category
// @Description To get menus from a category
// @Param id path string true "category id"
// @Accept json
// @Produce json
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /category/{id}/menus [get]
func (h *Handler) GetMenus(c *gin.Context) {
	categoryID := c.Param("id")

	menus, err := h.service.GetMenus(categoryID)
	if err != nil {
		shared.LogError("error getting menus from category", LogHandler, "GetMenus", err, menus)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingMenus))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(menus))
}
