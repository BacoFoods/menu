package assets

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{
		handler: handler,
	}
}

func (r *Routes) RegisterRoutes(routes *shared.CustomRoutes) {
	routes.GET("/assets/:code", r.handler.GetByPlaca)
	routes.POST("/assets", r.handler.CreateAsset)
}
