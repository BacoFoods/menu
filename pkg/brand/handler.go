package brand

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler = "pkg/brand/handler"

// Handler to handle requests to the brand service
type Handler struct {
	service Service
}

// NewHandler to create a new handler for the brand service
func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Find to handle a request to find all brands
// @Tags Brand
// @Summary To find brands
// @Description To find brands
// @Param name query string false "brand name"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Brand}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /brand [get]
func (h Handler) Find(c *gin.Context) {
	query := make(map[string]string)
	name := c.Query("name")
	if name != "" {
		query["name"] = c.Query("name")
	}
	brands, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding brands", LogHandler, "Find", err, brands)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingBrand))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(brands))
}

// Get to handle a request to get a brand
// @Tags Brand
// @Summary To get a brand
// @Description To get a brand
// @Param id path string true "brand id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Brand}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /brand/{id} [get]
func (h Handler) Get(c *gin.Context) {
	brandID := c.Param("id")
	brand, err := h.service.Get(brandID)
	if err != nil {
		shared.LogError("error getting brand", LogHandler, "Get", err, brand)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingBrand))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(brand))
}

// Create to handle a request to create a brand
// @Tags Brand
// @Summary To create a brand
// @Description To create a brand
// @Accept json
// @Produce json
// @Param brand body object{description=string,name=string} true "brand"
// @Success 200 {object} object{status=string,data=Brand}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /brand [post]
func (h Handler) Create(c *gin.Context) {
	var request Brand
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Create", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	brand, err := h.service.Create(&request)
	if err != nil {
		shared.LogError("error creating brand", LogHandler, "Create", err, brand)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingBrand))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(brand))
}

// Update to handle a request to update a brand
// @Tags Brand
// @Summary To update a brand
// @Description To update a brand
// @Accept json
// @Produce json
// @Param id path string true "brand id"
// @Param brand body Brand true "brand"
// @Success 200 {object} object{status=string,data=Brand}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /brand/{id} [patch]
func (h Handler) Update(c *gin.Context) {
	var request Brand
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Update", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	brand, err := h.service.Update(&request)
	if err != nil {
		shared.LogError("error updating brand", LogHandler, "Update", err, brand)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingBrand))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(brand))
}

// Delete to handle a request to delete a brand
// @Tags Brand
// @Summary To delete a brand
// @Description To delete a brand
// @Accept json
// @Produce json
// @Param id path string true "brand id"
// @Success 200 {object} object{status=string,data=Brand}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /brand/{id} [delete]
func (h Handler) Delete(c *gin.Context) {
	brandID := c.Param("id")
	brand, err := h.service.Delete(brandID)
	if err != nil {
		shared.LogError("error deleting brand", LogHandler, "Delete", err, brand)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingBrand))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(brand))
}
