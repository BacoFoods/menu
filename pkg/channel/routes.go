package channel

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/channel", r.handler.Find)
	router.GET("/channel/:id", r.handler.Get)
	router.POST("/channel", r.handler.Create)
	router.PATCH("/channel/:id", r.handler.Update)
	router.DELETE("/channel/:id", r.handler.Delete)
}
