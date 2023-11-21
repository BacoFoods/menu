package siesa

import (
	"github.com/gin-gonic/gin"
)

type Routes struct {
	handler *Handler
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.POST("/siesa", r.handler.Create)
	private.GET("/siesa/reference", r.handler.FindReferences)
	private.POST("/siesa/reference", r.handler.CreateReference)
	private.DELETE("/siesa/reference/:id", r.handler.DeleteReference)
	private.PATCH("/siesa/reference/:id", r.handler.UpdateReference)

}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}
