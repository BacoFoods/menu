package product

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/product/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Find to handle a request to find all products
// @Tags Product
// @Summary To find products
// @Description To find products
// @Param name query string false "product name"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /product [get]
func (h *Handler) Find(c *gin.Context) {
	query := make(map[string]string)
	name := c.Query("name")
	if name != "" {
		query["name"] = c.Query("name")
	}
	products, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding products", LogHandler, "Find", err, products)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingProduct))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(products))
}

// Get to handle a request to get a product
// @Tags Product
// @Summary To get a product
// @Description To get a product
// @Param id path string true "product id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /product/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	productID := c.Param("id")
	product, err := h.service.Get(productID)
	if err != nil {
		shared.LogError("error getting product", LogHandler, "Get", err, product)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingProduct))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(product))
}

// Create to handle a request to create a product
// @Tags Product
// @Summary To create a product
// @Description To create a product
// @Accept json
// @Produce json
// @Param product body Product true "product"
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /product [post]
func (h *Handler) Create(c *gin.Context) {
	var requestBody Product
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		shared.LogError("error binding request body", LogHandler, "Create", err, requestBody)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingProduct))
		return
	}

	product, err := h.service.Create(&requestBody)
	if err != nil {
		shared.LogError("error creating product", LogHandler, "Create", err, product)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingProduct))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(product))
}

// Update to handle a request to update a product
// @Tags Product
// @Summary To update a product
// @Description To update a product
// @Accept json
// @Produce json
// @Param id path string true "product id"
// @Param product body Product true "product"
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /product/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	var requestBody Product
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		shared.LogWarn("warning binding request body", LogHandler, "Update", err, requestBody)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}
	product, err := h.service.Update(&requestBody)
	if err != nil {
		shared.LogError("error updating product", LogHandler, "Update", err, product)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingProduct))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(product))
}

// Delete to handle a request to delete a product
// @Tags Product
// @Summary To delete a product
// @Description To delete a product
// @Accept json
// @Produce json
// @Param id path string true "product id"
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /product/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	productID := c.Param("id")
	product, err := h.service.Delete(productID)
	if err != nil {
		shared.LogError("error deleting product", LogHandler, "Delete", err, product)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingProduct))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(product))
}
