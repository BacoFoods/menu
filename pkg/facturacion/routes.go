package facturacion

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.POST("/store/:id/facturacion/config", r.handler.CreateConfig)
	private.GET("/store/:id/facturacion/config", r.handler.FindConfig)
	private.PUT("/store/:id/facturacion/config/:configId", r.handler.UpdateConfig)
}
