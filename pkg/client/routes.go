package client

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r *Routes) RegisterRoutes(routes *shared.CustomRoutes) {
	routes.GET("/clients", r.handler.List)
	routes.GET("/clients/:id", r.handler.Get)
	routes.POST("/clients", r.handler.Create)
	routes.PATCH("/clients/:id", r.handler.Update)
	routes.DELETE("/clients/:id", r.handler.Delete)
}
