package country

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *shared.CustomRoutes) {
	private.POST("/country", r.handler.Create)
	private.GET("/country", r.handler.Find)
	private.GET("/country/:id", r.handler.Get)
	private.PATCH("/country", r.handler.Update)
	private.DELETE("/country/:id", r.handler.Delete)
}
