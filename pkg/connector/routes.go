package connector

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.POST("/equivalence", r.handler.Create)
	private.POST("/connector", r.handler.CreateFile)
}
