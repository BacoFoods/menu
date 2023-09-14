package course

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/course/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Find to handle a request to find all course
// @Tags Course
// @Summary To find course
// @Description To find course
// @Param storeID query string false "storeID"
// @Param channelID query string false "channelID"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Course}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /course [get]
func (h *Handler) Find(ctx *gin.Context) {
	filter := make(map[string]any)

	if storeID := ctx.Query("storeID"); storeID != "" {
		filter["store_id"] = storeID
	}

	if channelID := ctx.Query("channelID"); channelID != "" {
		filter["channel_id"] = channelID
	}

	courses, err := h.service.Find(filter)
	if err != nil {
		shared.LogError("error finding courses", LogHandler, "Find", err, courses)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCourseFinding))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(courses))
}

// FindByID to handle a request to find a course by id
// @Tags Course
// @Summary To find a course by id
// @Description To find a course by id
// @Param id path string true "course id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Course}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /course/{id} [get]
func (h *Handler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	course, err := h.service.Get(id)
	if err != nil {
		shared.LogError("error finding course", LogHandler, "FindByID", err, course)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCourseGetting))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(course))
}

// Create to handle a request to create or update a course
// @Tags Course
// @Summary To create or update a course
// @Description Creates a course or updates if ID course is provided in body
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param course body Course true "course"
// @Success 200 {object} object{status=string,data=Course}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /course [put]
func (h *Handler) Create(ctx *gin.Context) {
	var course Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		shared.LogError("error binding course", LogHandler, "Create", err, course)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorCourseBinding))
		return
	}

	newCourse, err := h.service.Create(course)
	if err != nil {
		shared.LogError("error creating course", LogHandler, "Create", err, course)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCourseCreating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(newCourse))
}

// Delete to handle a request to delete a course
// @Tags Course
// @Summary To delete a course
// @Description To delete a course
// @Param id path string true "course id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Course}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /course/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	course, err := h.service.Delete(id)
	if err != nil {
		shared.LogError("error deleting course", LogHandler, "Delete", err, course)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCourseDeleting))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(course))
}
