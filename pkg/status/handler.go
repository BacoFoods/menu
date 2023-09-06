package status

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/status/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create
// @Tags Status
// @Summary To create a status
// @Description To create a status
// @Param status body CreateStatus true "status request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Status}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /status [post]
func (h *Handler) Create(ctx *gin.Context) {
	var dto CreateStatus
	if err := ctx.BindJSON(&dto); err != nil {
		shared.LogWarn("warning binding request fail", LogHandler, "Create", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	status, err := h.service.Create(dto.ToStatus())
	if err != nil {
		shared.LogError("error creating status", LogHandler, "Create", err, status)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorCreatingStatus))
		return
	}
	ctx.JSON(http.StatusOK, shared.SuccessResponse(*status))
}

// Delete
// @Tags Status
// @Summary To delete a status
// @Description To delete a status
// @Param id path string true "statusID"
// @Param status body Status true "status request"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=string}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /status/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.service.Delete(id); err != nil {
		shared.LogError("error deleting status", LogHandler, "Delete", err, id)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorDeletingStatus))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse("status deleted"))
}

// Find
// @Tags Status
// @Summary To find status
// @Description To find status
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Status}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /status [get]
func (h *Handler) Find(ctx *gin.Context) {
	statuses, err := h.service.Find()
	if err != nil {
		shared.LogError("error getting all status", LogHandler, "Find", err)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorGettingStatus))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(statuses))
}

// Get
// @Tags Status
// @Summary To find status by id
// @Description To find status by id
// @Param id path string true "status id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Status}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /status/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	statusID := ctx.Param("id")

	status, err := h.service.Get(statusID)
	if err != nil {
		shared.LogError("error getting status", LogHandler, "Get", err, statusID)
		ctx.JSON(http.StatusUnprocessableEntity, shared.SuccessResponse(ErrorGettingStatus))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(status))
}

// Update
// @Tags Status
// @Summary To update status
// @Description To update status
// @Param status body UpdateStatus true "status to update"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Status}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 403 {object} shared.Response
// @Router /status [patch]
func (h *Handler) Update(ctx *gin.Context) {
	var dto UpdateStatus
	if err := ctx.BindJSON(&dto); err != nil {
		shared.LogError("error getting status request body", LogHandler, "Update", err)
		ctx.JSON(http.StatusBadRequest, shared.ErrorResponse(ErrorBadRequest))
		return
	}

	status, err := h.service.Update(dto.ToStatus())
	if err != nil {
		shared.LogError("error updating status", LogHandler, "Update", err, status)
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorUpdatingStatus))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(status))
}
