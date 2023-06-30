package menu

import (
	_ "github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service,
	}
}

// Create to handle request for create a menu
// @Tags Menu
// @Summary To create a menu
// @Description To create a menu
// @Param menu body Menu true "menu request"
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /menu [post]
func (h *Handler) Create(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "create menu"})
}

// Find to handle request for find menus
// @Tags Menu
// @Summary To find menu
// @Description To find menu
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /menu [get]
func (h *Handler) Find(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "find menu"})
}

// Get to handle request for get a menu
// @Tags Menu
// @Summary To get a menu
// @Description To get a menu
// @Accept json
// @Produce json
// @Param id path string true "menu id"
// @Success 200 {object} object{message=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /menu/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "get menu"})
}

// Update to handle request for update a menu
// @Tags Menu
// @Summary To update a menu
// @Description To update a menu
// @Param menu body Menu true "menu request"
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /menu [patch]
func (h *Handler) Update(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "update menu"})
}

// Delete to handle request for delete a menu
// @Tags Menu
// @Summary To delete a menu
// @Description To delete a menu
// @Param id path string true "menu id"
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string,data=Menu}
// @Failure 400 {object} shared.Response
// @Failure 422 {object} shared.Response
// @Router /menu/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "delete menu"})
}
