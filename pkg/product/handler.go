package product

import (
	"net/http"
	"strconv"

	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

const LogHandler string = "pkg/product/handler"

type RequestUpdateOverriders struct {
	Field string `json:"field" example:"price"`
	Value string `json:"value" example:"10000"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Product

// Find to handle a request to find all products
// @Tags Product
// @Summary To find products
// @Description To find products
// @Param name query string false "product name"
// @Param brand-id query string false "brand id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /product [get]
func (h *Handler) Find(c *gin.Context) {
	query := make(map[string]string)

	name := c.Query("name")
	if name != "" {
		query["name"] = c.Query("name")
	}

	brandID := c.Query("brand-id")
	if brandID != "" {
		query["brand_id"] = brandID
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
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /product/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	productID := c.Param("id")

	products, err := h.service.Get(productID)
	if err != nil {
		shared.LogError("error getting product", LogHandler, "Get", err, products)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingProduct))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(products))
}

// Create to handle a request to create a product
// @Tags Product
// @Summary To create a product
// @Description This endpoint creates a product and makes an association with all store channels to this product to allow
//
//	the product to be available in all channels by default with ove
//
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param product body Product true "product"
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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
// @Security ApiKeyAuth
// @Param id path string true "product id"
// @Param product body Product true "product"
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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
// @Security ApiKeyAuth
// @Param id path string true "product id"
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
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

// AddModifier to handle a request to add a modifier to a product
// @Tags Product
// @Summary To add a modifier to a product
// @Description To add a modifier to a product
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "product id"
// @Param modifierID path string true "modifier id"
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /product/{id}/modifier/{modifierID} [post]
func (h *Handler) AddModifier(c *gin.Context) {
	productID := c.Param("id")
	modifierID := c.Param("modifierID")

	product, err := h.service.AddModifier(productID, modifierID)
	if err != nil {
		shared.LogError("error adding modifier to product", LogHandler, "AddModifier", err, product)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorAddingModifier))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(product))
}

// RemoveModifier to handle a request to remove a modifier from a product
// @Tags Product
// @Summary To remove a modifier from a product
// @Description To remove a modifier from a product
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "product id"
// @Param modifierID path string true "modifier id"
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /product/{id}/modifier/{modifierID} [delete]
func (h *Handler) RemoveModifier(c *gin.Context) {
	productID := c.Param("id")
	modifierID := c.Param("modifierID")

	product, err := h.service.RemoveModifier(productID, modifierID)
	if err != nil {
		shared.LogError("error removing modifier from product", LogHandler, "RemoveModifier", err, product)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorRemovingModifier))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(product))
}

// GetOverridersByField to handle a request to get overriders for a product
// @Tags Product
// @Summary To get overriders for a product
// @Description To get overriders for a product
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "product id"
// @Param field query string false "options: (price, name, description, image, enable)"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /product/{id}/overrider [get]
func (h *Handler) GetOverridersByField(c *gin.Context) {
	productID := c.Param("id")

	field, ok := Entities[c.Query("field")]
	if !ok {
		shared.LogWarn("warning getting overriders for product", LogHandler, "GetOverridersByField", nil, productID, c.Query("field"))
		c.JSON(http.StatusOK, shared.SuccessResponse([]Overrider{}))
		return
	}

	overriders, err := h.service.GetOverriders(productID, field.Code)
	if err != nil {
		shared.LogError("error getting overriders for product", LogHandler, "GetOverriders", err, productID, field)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse("Error getting overriders for product"))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(overriders))
}

// UpdateAllOverriders to handle a request to update all overriders for a product
// @Tags Product
// @Summary To update all overriders for a product
// @Description To update all overriders for a product
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "product id"
// @Param request body RequestUpdateOverriders true "request"
// @Success 200 {object} object{status=string,data=Product}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /product/{id}/overrider/update-all [patch]
func (h *Handler) UpdateAllOverriders(c *gin.Context) {
	productID := c.Param("id")

	var request RequestUpdateOverriders
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogError("error updating all overriders for product", LogHandler, "UpdateAllOverriders", err, productID, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse("Error updating all overriders for product"))
		return
	}

	if _, ok := Entities[request.Field]; !ok {
		shared.LogWarn("warning updating all overriders for product", LogHandler, "UpdateAllOverriders", nil, productID, request)
		c.JSON(http.StatusOK, shared.SuccessResponse(ErrorBadRequest))
		return
	}

	value := TransformValue(request.Field, request.Value)

	err := h.service.UpdateAllOverriders(productID, request.Field, value)
	if err != nil {
		shared.LogError("error updating all overriders for product", LogHandler, "UpdateAllOverriders", err, productID, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse("Error updating all overriders for product"))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(err))
}

// GetCategories to handle a request to get categories for a product
// @Tags Product
// @Summary To get categories for a product
// @Description To get categories for a product
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "product id"
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /product/{id}/category [get]
func (h *Handler) GetCategories(c *gin.Context) {
	productID := c.Param("id")

	categories, err := h.service.GetCategory(productID)
	if err != nil {
		shared.LogError("error getting categories for product", LogHandler, "GetCategories", err, productID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingCategory))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(categories))
}

// Modifiers

// ModifierFind to handle a request to find modifiers
// @Tags Modifiers
// @Summary To find modifiers
// @Description To find modifiers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "modifier name"
// @Param brand-id query string false "brand id"
// @Success 200 {object} object{status=string,data=Modifier}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /modifier [get]
func (h *Handler) ModifierFind(c *gin.Context) {
	query := make(map[string]string)

	name := c.Query("name")
	if name != "" {
		query["name"] = name
	}

	brandID := c.Query("brand-id")
	if brandID != "" {
		query["brand_id"] = brandID
	}

	modifiers, err := h.service.ModifierFind(query)
	if err != nil {
		shared.LogError("error finding modifiers", LogHandler, "FindModifiers", err, modifiers)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorModifierGetting))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(modifiers))
}

// ModifierCreate to handle a request to create a modifier
// @Tags Modifiers
// @Summary To create a modifier
// @Description To create a modifier
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param modifier body Modifier true "modifier"
// @Success 200 {object} object{status=string,data=Modifier}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /modifier [post]
func (h *Handler) ModifierCreate(c *gin.Context) {
	var body Modifier
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "CreateModifier", err, body)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	modifier, err := h.service.ModifierCreate(&body)
	if err != nil {
		shared.LogError("error creating modifier", LogHandler, "CreateModifier", err, modifier)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorModifierCreation))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(modifier))
}

// ModifierAddProduct to handle a request to add a product to a modifier
// @Tags Modifiers
// @Summary To add a product to a modifier
// @Description To add a product to a modifier
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "modifier id"
// @Param productID path string true "product id"
// @Success 200 {object} object{status=string,data=Modifier}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /modifier/{id}/product/{productID} [post]
func (h *Handler) ModifierAddProduct(c *gin.Context) {
	modifierID := c.Param("id")
	productID := c.Param("productID")

	modifier, err := h.service.ModifierAddProduct(productID, modifierID)
	if err != nil {
		shared.LogError("error adding product to modifier", LogHandler, "ModifierAddProduct", err, modifier)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorModifierAddingProduct))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(modifier))
}

// ModifierRemoveProduct to handle a request to remove a product from a modifier
// @Tags Modifiers
// @Summary To remove a product from a modifier
// @Description To remove a product from a modifier
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "modifier id"
// @Param productID path string true "product id"
// @Success 200 {object} object{status=string,data=Modifier}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /modifier/{id}/product/{productID} [delete]
func (h *Handler) ModifierRemoveProduct(c *gin.Context) {
	modifierID := c.Param("id")
	productID := c.Param("productID")

	modifier, err := h.service.ModifierRemoveProduct(productID, modifierID)
	if err != nil {
		shared.LogError("error removing product from modifier", LogHandler, "ModifierRemoveProduct", err, modifier)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorModifierRemovingProduct))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(modifier))
}

// ModifierUpdate to handle a request to update a modifier
// @Tags Modifiers
// @Summary To update a modifier
// @Description To update a modifier
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "modifier id"
// @Param modifier body ModifierDTO true "modifier"
// @Success 200 {object} object{status=string,data=Modifier}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /modifier/{id} [patch]
func (h *Handler) ModifierUpdate(c *gin.Context) {
	modifierID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		shared.LogError("error parsing modifier id", LogHandler, "ModifierUpdate", err, modifierID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	var body ModifierDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogError("error binding request body", LogHandler, "ModifierUpdate", err, body)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	modifier := body.ToModifier()
	modifier.ID = uint(modifierID)
	modifierUpdated, err := h.service.ModifierUpdate(&modifier)
	if err != nil {
		shared.LogError("error updating modifier", LogHandler, "ModifierUpdate", err, modifier)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorModifierUpdate))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(modifierUpdated))
}

// Overriders

// OverriderFind to handle a request to find all overriders
// @Tags Overrider
// @Summary To find overrider
// @Description To find overrider
// @Param name query string false "overrider name"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Overrider}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /overrider [get]
func (h *Handler) OverriderFind(c *gin.Context) {
	query := make(map[string]string)
	name := c.Query("name")
	if name != "" {
		query["name"] = c.Query("name")
	}

	overriders, err := h.service.OverriderFind(query)
	if err != nil {
		shared.LogError("error finding overriders", LogHandler, "OverriderFind", err, overriders)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOverriderFinding))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(overriders))
}

// OverriderGet to handle a request to get an overrider
// @Tags Overrider
// @Summary To get an overrider
// @Description To get an overrider
// @Param id path string true "overrider id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Overrider}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /overrider/{id} [get]
func (h *Handler) OverriderGet(c *gin.Context) {
	overriderID := c.Param("id")

	overrider, err := h.service.OverriderGet(overriderID)
	if err != nil {
		shared.LogError("error getting overrider", LogHandler, "OverriderGet", err, overrider)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOverriderGetting))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(overrider))
}

// OverriderCreate to handle a request to create an overrider
// @Tags Overrider
// @Summary To create an overrider
// @Description To create an overrider
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param overriders body Overrider true "overrider"
// @Success 200 {object} object{status=string,data=Overrider}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /overrider [post]
func (h *Handler) OverriderCreate(c *gin.Context) {
	var request Overrider
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding overrider", LogHandler, "OverriderCreate", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	overrider, err := h.service.OverriderCreate(&request)
	if err != nil {
		shared.LogError("error creating overrider", LogHandler, "OverriderCreate", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOverriderCreating))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(overrider))
}

// OverriderUpdate to handle a request to update an overrider
// @Tags Overrider
// @Summary To update an overrider
// @Description To update an overrider
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "overrider id"
// @Param overriders body Overrider true "overrider"
// @Success 200 {object} object{status=string,data=Overrider}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /overrider/{id} [patch]
func (h *Handler) OverriderUpdate(c *gin.Context) {
	var request Overrider
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding overrider", LogHandler, "OverriderUpdate", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	overrider, err := h.service.OverriderUpdate(&request)
	if err != nil {
		shared.LogError("error updating overrider", LogHandler, "Update", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOverriderUpdating))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(overrider))
}

// OverriderDelete to handle a request to delete an overrider
// @Tags Overrider
// @Summary To delete an overrider
// @Description To delete an overrider
// @Param id path string true "overrider id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Overrider}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /overrider/{id} [delete]
func (h *Handler) OverriderDelete(c *gin.Context) {
	overriderID := c.Param("id")

	overrider, err := h.service.Delete(overriderID)
	if err != nil {
		shared.LogError("error deleting overrider", LogHandler, "OverriderDelete", err, overrider)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorOverriderDeleting))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(overrider))
}
