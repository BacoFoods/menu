package discount

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.POST("/discount", r.handler.Create)
	private.GET("/discount", r.handler.Find)
	private.GET("/discount/:id", r.handler.Get)
	private.PATCH("/discount", r.handler.Update)
	private.DELETE("/discount/:id", r.handler.Delete)
}
