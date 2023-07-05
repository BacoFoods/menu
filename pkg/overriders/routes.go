package overriders

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/overriders", r.handler.Find)
	router.GET("/overriders/:id", r.handler.Get)
	router.POST("/overriders", r.handler.Create)
	router.PATCH("/overriders/:id", r.handler.Update)
	router.DELETE("/overriders/:id", r.handler.Delete)
}
