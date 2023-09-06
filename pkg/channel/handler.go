package channel

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/channel/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Find to handle a request to find all channels
// @Tags Channel
// @Summary To find channels
// @Description To find channels
// @Param name query string false "channel name"
// @Param store-id query string false "store id"
// @Param brand-id query string false "brand id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Channel}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /channel [get]
func (h *Handler) Find(c *gin.Context) {
	query := make(map[string]string)

	name := c.Query("name")
	if name != "" {
		query["name"] = name
	}

	storeID := c.Query("store-id")
	if storeID != "" {
		query["store_id"] = storeID
	}

	brandID := c.Query("brand-id")
	if brandID != "" {
		query["brand_id"] = brandID
	}

	channels, err := h.service.Find(query)
	if err != nil {
		shared.LogError("error finding channels", LogHandler, "Find", err, channels)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorFindingChannel))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(channels))
}

// Get to handle a request to get a channel
// @Tags Channel
// @Summary To get a channel
// @Description To get a channel
// @Param id path string true "channel id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Channel}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /channel/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	channelID := c.Param("id")
	channel, err := h.service.Get(channelID)
	if err != nil {
		shared.LogError("error getting channel", LogHandler, "Get", err, channel)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingChannel))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(channel))
}

// Create to handle a request to create a channel
// @Tags Channel
// @Summary To create a channel
// @Description To create a channel
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param channel body Channel true "channel"
// @Success 200 {object} object{status=string,data=Channel}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /channel [post]
func (h *Handler) Create(c *gin.Context) {
	var request Channel
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Create", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorBadRequest))
		return
	}
	channel, err := h.service.Create(&request)
	if err != nil {
		shared.LogError("error creating channel", LogHandler, "Create", err, channel)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingChannel))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(channel))

}

// Update to handle a request to update a channel
// @Tags Channel
// @Summary To update a channel
// @Description To update a channel
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "channel id"
// @Param channel body Channel true "channel"
// @Success 200 {object} object{status=string,data=Channel}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /channel/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	var request Channel
	if err := c.ShouldBindJSON(&request); err != nil {
		shared.LogWarn("warning binding request", LogHandler, "Update", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	channel, err := h.service.Update(&request)
	if err != nil {
		shared.LogError("error updating channel", LogHandler, "Update", err, request)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingChannel))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(channel))
}

// Delete to handle a request to delete a channel
// @Tags Channel
// @Summary To delete a channel
// @Description To delete a channel
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "channel id"
// @Success 200 {object} object{status=string,data=Channel}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /channel/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	channelID := c.Param("id")
	channel, err := h.service.Delete(channelID)
	if err != nil {
		shared.LogError("error deleting channel", LogHandler, "Delete", err, channel)
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingChannel))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(channel))
}
