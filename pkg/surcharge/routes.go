package surcharge

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *shared.CustomRoutes) {
	router.GET("/surcharge", r.handler.Find)
	router.GET("/surcharge/:id", r.handler.Get)
	router.POST("/surcharge", r.handler.Create)
	router.PATCH("/surcharge/:id", r.handler.Update)
	router.DELETE("/surcharge/:id", r.handler.Delete)
}
