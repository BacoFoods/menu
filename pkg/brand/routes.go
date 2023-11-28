package brand

import (
	"github.com/BacoFoods/menu/pkg/shared"
)

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r *Routes) RegisterRoutes(router *shared.CustomRoutes) {
	router.GET("/brand", r.handler.Find)
	router.GET("/brand/:id", r.handler.Get)
	router.POST("/brand", r.handler.Create)
	router.PATCH("/brand/:id", r.handler.Update)
	router.DELETE("/brand/:id", r.handler.Delete)
}
