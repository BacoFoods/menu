package taxes

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.POST("/tax", r.handler.Create)
	private.GET("/tax", r.handler.Find)
	private.GET("/tax/:id", r.handler.Get)
	private.PATCH("/tax", r.handler.Update)
	private.DELETE("/tax/:id", r.handler.Delete)
}
