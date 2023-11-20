package category

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
// @Param name query string false "category by name"
// @Param brandID query string false "category by brand"
// @Param productID query string false "category by product"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /category [get]
func (h *Handler) Find(c *gin.Context) {
	query := make(map[string]string)

	name := c.Query("name")
	if name != "" {
		query["name"] = c.Query("name")
	}

	brandID := c.Query("brandID")
	if brandID != "" {
		query["brand_id"] = brandID
	}

	productID := c.Query("productID")
	if productID != "" {
		query["product_id"] = productID
	}

	categories, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding categories", LogHandler, "Find", err, categories)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCategoryFinding))
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
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /category/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	categoryID := c.Param("id")
	category, err := h.service.Get(categoryID)
	if err != nil {
		shared.LogError("error getting category", LogHandler, "Get", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCategoryGetting))
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
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /category [post]
func (h *Handler) Create(c *gin.Context) {
	var requestBody Category
	if err := c.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCategoryBadRequest))
		return
	}

	category, err := h.service.Create(&requestBody)
	if err != nil {
		shared.LogError("error creating category", LogHandler, "Create", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCategoryCreating))
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
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /category/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	var requestBody Category
	if err := c.BindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Update", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCategoryBadRequest))
		return
	}
	category, err := h.service.Update(&requestBody)
	if err != nil {
		shared.LogError("error updating category", LogHandler, "Update", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCategoryUpdating))
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
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Category}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /category/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	categoryID := c.Param("id")
	category, err := h.service.Delete(categoryID)
	if err != nil {
		shared.LogError("error deleting category", LogHandler, "Delete", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCategoryDeleting))
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
// @Security ApiKeyAuth
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /category/{id}/menu [get]
func (h *Handler) GetMenus(c *gin.Context) {
	categoryID := c.Param("id")

	menus, err := h.service.GetMenus(categoryID)
	if err != nil {
		shared.LogError("error getting menus from category", LogHandler, "GetMenus", err, menus)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCategoryGettingMenus))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(menus))
}

// AddProduct to handle a request to add a product to a category
// @Tags Category
// @Summary To add a product to a category
// @Description To add a product to a category
// @Param id path string true "category id"
// @Param productID path string true "product id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /category/{id}/product/{productID}/add [patch]
func (h *Handler) AddProduct(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing category id fail", LogHandler, "AddProduct", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCategoryBadRequest))
		return
	}

	productID, err := strconv.ParseUint(c.Param("productID"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing product id fail", LogHandler, "AddProduct", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCategoryBadRequest))
		return
	}

	category, err := h.service.AddProduct(uint(categoryID), uint(productID))
	if err != nil {
		shared.LogError("error adding product to category", LogHandler, "AddProduct", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCategoryAddingProduct))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(category))
}

// RemoveProduct to handle a request to remove a product from a category
// @Tags Category
// @Summary To remove a product from a category
// @Description To remove a product from a category
// @Param id path integer true "category id"
// @Param productID path integer true "product id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /category/{id}/product/{productID}/remove [patch]
func (h *Handler) RemoveProduct(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing category id fail", LogHandler, "RemoveProduct", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCategoryBadRequest))
		return
	}

	productID, err := strconv.ParseUint(c.Param("productID"), 10, 64)
	if err != nil {
		shared.LogWarn("warning parsing product id fail", LogHandler, "RemoveProduct", err)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCategoryBadRequest))
		return
	}

	category, err := h.service.RemoveProduct(uint(categoryID), uint(productID))
	if err != nil {
		shared.LogError("error removing product from category", LogHandler, "RemoveProduct", err, category)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCategoryRemovingProduct))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(category))
}
