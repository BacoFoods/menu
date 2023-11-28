package temporal

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
	routes.GET("/temporal/locales", r.handler.GetLocales)
	routes.GET("/temporal/arqueo", r.handler.GetArqueo)
}
