package shift

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/shift-cashier/open", r.handler.OpenShift)
	router.POST("/shift-cashier/close", r.handler.CloseShift)
}
