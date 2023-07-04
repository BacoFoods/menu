package spot

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler = "pkg/spot/handler"

// Handler to handle requests to the spot service
type Handler struct {
	service Service
}

// NewHandler to create a new handler for the spot service
func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Find to handle a request to find all spots
// @Tags Spot
// @Summary To find spots
// @Description To find spots
// @Param name query string false "spot name"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=[]Spot}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /spot [get]
func (h Handler) Find(c *gin.Context) {
	query := make(map[string]string)
	name := c.Query("name")
	if name != "" {
		query["name"] = c.Query("name")
	}
	spots, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding spots", LogHandler, "Find", err, spots)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingSpot))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(spots))
}

// Get to handle a request to get a spot
// @Tags Spot
// @Summary To get a spot
// @Description To get a spot
// @Param id path string true "spot id"
// @Accept json
// @Produce json
// @Success 200 {object} object{status=string,data=Spot}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /spot/{id} [get]
func (h Handler) Get(c *gin.Context) {
	spotID := c.Param("id")
	spot, err := h.service.Get(spotID)
	if err != nil {
		shared.LogError("error getting spot", LogHandler, "Get", err, spotID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingSpot))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(spot))
}

// Create to handle a request to create a spot
// @Tags Spot
// @Summary To create a spot
// @Description To create a spot
// @Accept json
// @Produce json
// @Param spot body Spot true "spot"
// @Success 200 {object} object{status=string,data=Spot}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /spot [post]
func (h Handler) Create(c *gin.Context) {
	var request Spot
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogError("error binding json", LogHandler, "Create", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	spot, err := h.service.Create(&request)
	if err != nil {
		shared.LogError("error creating spot", LogHandler, "Create", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingSpot))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(spot))
}

// Update to handle a request to update a spot
// @Tags Spot
// @Summary To update a spot
// @Description To update a spot
// @Accept json
// @Produce json
// @Param id path string true "spot id"
// @Param spot body Spot true "spot"
// @Success 200 {object} object{status=string,data=Spot}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /spot/{id} [patch]
func (h Handler) Update(c *gin.Context) {
	var request Spot
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Update", err, request)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	spot, err := h.service.Update(&request)
	if err != nil {
		shared.LogError("error updating spot", LogHandler, "Update", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingSpot))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(spot))
}

// Delete to handle a request to delete a spot
// @Tags Spot
// @Summary To delete a spot
// @Description To delete a spot
// @Accept json
// @Produce json
// @Param id path string true "spot id"
// @Success 200 {object} object{status=string,data=Spot}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /spot/{id} [delete]
func (h Handler) Delete(c *gin.Context) {
	spotID := c.Param("id")
	spot, err := h.service.Delete(spotID)
	if err != nil {
		shared.LogError("error deleting spot", LogHandler, "Delete", err, spotID)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingSpot))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(spot))
}
