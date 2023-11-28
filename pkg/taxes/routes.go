package taxes

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *shared.CustomRoutes) {
	private.POST("/tax", r.handler.Create)
	private.GET("/tax", r.handler.Find)
	private.GET("/tax/:id", r.handler.Get)
	private.PATCH("/tax", r.handler.Update)
	private.DELETE("/tax/:id", r.handler.Delete)
}
