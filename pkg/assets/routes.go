package assets

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{
		handler: handler,
	}
}

func (r *Routes) RegisterRoutes(routes *gin.RouterGroup) {
	routes.GET("/assets/:code", r.handler.GetByPlaca)
	routes.POST("/assets", r.handler.CreateAsset)
}