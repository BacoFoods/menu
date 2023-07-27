package overriders

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.GET("/overriders", r.handler.Find)
	private.GET("/overriders/:id", r.handler.Get)
	private.POST("/overriders", r.handler.Create)
	private.PATCH("/overriders/:id", r.handler.Update)
	private.DELETE("/overriders/:id", r.handler.Delete)
}
