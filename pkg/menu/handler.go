package menu

import (
	availabilityPkg "github.com/BacoFoods/menu/pkg/availability"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/menu/handler"

type RequestMenuCreate struct {
	Name     string `json:"name" binding:"required"`
	BrandID  uint   `json:"brand_id" binding:"required"`
	Place    string `json:"place" binding:"required"`
	PlaceIDs []uint `json:"place_id" binding:"required"`
}

type RequestMenuAvailability struct {
	PlaceIDs map[uint]bool `json:"places" binding:"required"`
}

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
// @Param brand-id query string false "brand id"
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
		query["name"] = name
	}

	brandID := c.Query("brand-id")
	if brandID != "" {
		query["brand_id"] = brandID
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
// @Param menu body RequestMenuCreate true "menu"
// @Success 200 {object} object{status=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /menu [post]
func (h *Handler) Create(c *gin.Context) {
	var body RequestMenuCreate
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogWarn("warning binding request body", LogHandler, "Create", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	menu, err := h.service.Create(body.Name, body.BrandID, body.Place, body.PlaceIDs)
	if err != nil {
		shared.LogError("error creating menu", LogHandler, "Create", err, body)
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

// ListByPlace to handle a request to list menus by place
// @Tags Menu
// @Summary To list menus by place
// @Description To list menus by place
// @Param place path string true "place"
// @Param place-id path string true "place id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /menu/place/{place}/{place-id}/list [get]
func (h *Handler) ListByPlace(c *gin.Context) {
	place := c.Param("place")
	placeID := c.Param("place-id")
	menus, err := h.service.FindByPlace(place, placeID)
	if err != nil {
		shared.LogError("error finding menus", LogHandler, "FindByPlace", err, menus)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingByPlace))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(menus))
}

// GetByPlace to handle a request to get a menu by place and load overriders
// @Tags Menu
// @Summary To get a menu by place and load overriders
// @Description To get a menu by place and load overriders
// @Param place path string true "place"
// @Param place-id path string true "place id"
// @Param menu-id path string true "menu id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /menu/place/{place}/{place-id}/menu-id/{menu-id} [get]
func (h *Handler) GetByPlace(c *gin.Context) {
	place := c.Param("place")
	placeID := c.Param("place-id")
	menuID := c.Param("menu-id")

	menu, err := h.service.GetByPlace(place, placeID, menuID)

	if err != nil {
		shared.LogError("error getting menu by place", LogHandler, "GetByPlace", err, menu)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingMenu))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(menu))
}

// UpdateAvailability to handle a request to update availability of a menu
// @Tags Menu
// @Summary To update availability of a menu
// @Description To update availability of a menu
// @Param menu-id path string true "menu id"
// @Param place path string true "place"
// @Param availability body RequestMenuAvailability true "availability"
// @Accept json
// @Produce json
// @Success 200 {object} shared.Response
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /menu/:id/availability/ [put]
func (h *Handler) UpdateAvailability(c *gin.Context) {
	place, err := availabilityPkg.GetPlace(c.Param("place"))
	if err != nil {
		shared.LogWarn("warning getting place", LogHandler, "UpdateAvailability", err, place)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	var body RequestMenuAvailability
	if err := c.ShouldBindJSON(&body); err != nil {
		shared.LogWarn("warning binding request body", LogHandler, "UpdateAvailability", err, body)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	menuID := c.Param("id")
	if _, err := h.service.UpdateAvailability(menuID, string(place), body.PlaceIDs); err != nil {
		shared.LogError("error updating availability", LogHandler, "UpdateAvailability", err, body)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingAvailability))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(nil))
}
