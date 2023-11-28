package shift

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(router *shared.CustomRoutes) {
	router.POST("/shift/open", r.handler.Open)
	router.POST("/shift/close", r.handler.Close)
}
