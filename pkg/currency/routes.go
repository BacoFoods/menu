package currency

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *shared.CustomRoutes) {
	private.POST("/currency", r.handler.Create)
	private.GET("/currency", r.handler.Find)
	private.GET("/currency/:id", r.handler.Get)
	private.PATCH("/currency", r.handler.Update)
	private.DELETE("/currency/:id", r.handler.Delete)
}
