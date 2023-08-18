package status

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.POST("/status", r.handler.Create)
	private.DELETE("/status/:id", r.handler.Delete)
	private.GET("/status", r.handler.Find)
	private.GET("/status/:id", r.handler.Get)
	private.PATCH("/status", r.handler.Update)
}
