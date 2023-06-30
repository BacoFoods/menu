package country

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) Register(private *gin.RouterGroup) {
	private.POST("/country", r.handler.Create)
	private.GET("/country", r.handler.Find)
	private.GET("/country/:id", r.handler.Get)
	private.PATCH("/country", r.handler.Update)
	private.DELETE("/country/:id", r.handler.Delete)
}
