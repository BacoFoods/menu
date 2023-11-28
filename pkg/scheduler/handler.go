package scheduler

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Find to handle a request to find all schedules
// @Tags Schedule
// @Summary To find schedules
// @Description To find schedules
// @Param store_id query string false "store id"
// @Param brand_id query string false "brand id"
// @Param day query string false "day" Enums(monday, tuesday, wednesday, thursday, friday, saturday, sunday, holiday)
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=ResponseBrand}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /schedules [get]
func (h *Handler) Find(c *gin.Context) {
	filter := make(map[string]any)

	if c.Query("store_id") != "" {
		filter["store_id"] = c.Query("store_id")
	}

	if c.Query("brand_id") != "" {
		filter["brand_id"] = c.Query("brand_id")
	}

	if c.Query("day") != "" {
		filter["day"] = c.Query("day")
	}

	schedules, err := h.service.Find(filter)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	responseBrand := ResponseBrand{}
	stores := make(map[uint]*ResponseStore)
	for _, schedule := range schedules {
		store, ok := stores[schedule.Store.ID]
		if !ok {
			stores[schedule.Store.ID] = &ResponseStore{
				ID:   schedule.Store.ID,
				Name: schedule.Store.Name,
				Schedules: []ResponseSchedules{
					{
						ID:      schedule.ID,
						Day:     schedule.Day,
						StoreID: schedule.StoreID,
						BrandID: schedule.BrandID,
						Opening: schedule.Opening,
						Closing: schedule.Closing,
						Enable:  schedule.Enable,
					},
				},
			}
		} else {
			store.Schedules = append(stores[schedule.Store.ID].Schedules, ResponseSchedules{
				ID:      schedule.ID,
				Day:     schedule.Day,
				StoreID: schedule.StoreID,
				BrandID: schedule.BrandID,
				Opening: schedule.Opening,
				Closing: schedule.Closing,
				Enable:  schedule.Enable,
			})
		}

		if strings.ToLower(time.Now().Weekday().String()) == schedule.Day {
			isOpen := schedule.IsOpen()
			stores[schedule.Store.ID].Open = isOpen
		}
	}

	for _, store := range stores {
		responseBrand.Stores = append(responseBrand.Stores, *store)
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(responseBrand))
}

// Create to handle a request to create a schedule
// @Tags Schedule
// @Summary To create a schedule
// @Description To create a schedule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param schedule body RequestSchedule true "schedule"
// @Success 200 {object} object{status=string,data=Schedule}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /schedules [post]
func (h *Handler) Create(c *gin.Context) {
	var schedule RequestSchedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	if err := schedule.RequestValidate(); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	if err := h.service.Create(schedule.ToSchedule()); err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(schedule))
}

// Update to handle a request to update a schedule
// @Tags Schedule
// @Summary To update a schedule
// @Description To update a schedule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param schedule body RequestSchedule true "schedule"
// @Success 200 {object} object{status=string,data=Schedule}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /schedules [patch]
func (h *Handler) Update(c *gin.Context) {
	var schedule RequestSchedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	if err := schedule.RequestValidate(); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	if err := h.service.Update(schedule.ToSchedule()); err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(schedule))
}

// Delete to handle a request to delete a schedule
// @Tags Schedule
// @Summary To delete a schedule
// @Description To delete a schedule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param schedule body Schedule true "schedule"
// @Success 200 {object} object{status=string,data=Schedule}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /schedules [delete]
func (h *Handler) Delete(c *gin.Context) {
	var schedule Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	if err := h.service.Delete(&schedule); err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// TodayStore to handle a request to find today's schedule of a store
// @Tags Schedule
// @Summary To find today's schedule of a store
// @Description To find today's schedule of a store
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "store id"
// @Success 200 {object} object{status=string,data=Schedule}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /schedules/store/{id}/today [get]
func (h *Handler) TodayStore(c *gin.Context) {
	schedule, err := h.service.Today(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(schedule))
}

// TodayBrand to handle a request to find today's store's schedules of a brand
// @Tags Schedule
// @Summary To find today's store's schedules of a brand
// @Description To find today's store's schedules of a brand
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "brand id"
// @Success 200 {object} object{status=string,data=[]Schedule}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /schedules/brand/{id}/today [get]
func (h *Handler) TodayBrand(c *gin.Context) {
	schedules, err := h.service.TodayStores(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, shared.SuccessResponse(schedules))
}

// EnableStore to handle a request to enable a store's schedule
// @Tags Schedule
// @Summary Turns enable a store's schedule
// @Description To enable a store's schedule
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "store id"
// @Param schedule body RequestEnableStore true "schedule"
// @Success 200 {object} object{status=string,data=Schedule}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /schedules/store/{id}/enable [post]
func (h *Handler) EnableStore(c *gin.Context) {
	var request RequestEnableStore
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse(err.Error()))
		return
	}

	schedules, err := h.service.EnableStore(c.Param("id"), request.Enable)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessResponse(schedules))
}
