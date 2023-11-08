package shift

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/shift/open", r.handler.Open)
	router.POST("/shift/close", r.handler.Close)
}
