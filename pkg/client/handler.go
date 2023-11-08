package client

import (
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const LogHandler string = "pkg/client/handler"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// List to handle a request to list all clients
// @Tags Client
// @Summary To list all clients
// @Description To list all clients
// @Accept json
// @Produce json
// @Param document query string false "document"
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=[]Client}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /clients [get]
func (h *Handler) List(ctx *gin.Context) {
	filter := make(map[string]any)
	if document := ctx.Query("document"); document != "" {
		filter["document"] = document
	}
	clients, err := h.service.List(filter)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorClientListing))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(clients))
}

// Get to handle a request to get a client
// @Tags Client
// @Summary To get a client
// @Description To get a client
// @Param id path string true "client id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Client}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /clients/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	client, err := h.service.Get(id)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorClientGetting))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(client))
}

// Create to handle a request to create a client
// @Tags Client
// @Summary To create a client
// @Description To create a client
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param client body Client true "client"
// @Success 200 {object} object{status=string,data=Client}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /clients [post]
func (h *Handler) Create(ctx *gin.Context) {
	var body Client
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorClientCreating))
		return
	}

	client, err := h.service.Create(&body)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorClientCreating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(client))
}

// Update to handle a request to update a client
// @Tags Client
// @Summary To update a client
// @Description To update a client
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param client body Client true "client"
// @Success 200 {object} object{status=string,data=Client}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /clients/{id} [patch]
func (h *Handler) Update(ctx *gin.Context) {
	var body Client
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorClientUpdating))
		return
	}

	client, err := h.service.Update(&body)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorClientUpdating))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(client))
}

// Delete to handle a request to delete a client
// @Tags Client
// @Summary To delete a client
// @Description To delete a client
// @Param id path string true "client id"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} object{status=string,data=Client}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Failure 401 {object} shared.Response
// @Router /clients/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	client, err := h.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, shared.ErrorResponse(ErrorClientDeleting))
		return
	}

	ctx.JSON(http.StatusOK, shared.SuccessResponse(client))
}
