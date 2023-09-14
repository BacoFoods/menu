package course

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.GET("/course", r.handler.Find)
	private.GET("/course/:id", r.handler.FindByID)
	private.PUT("/course", r.handler.Create)
	private.DELETE("/course/:id", r.handler.Delete)
}
