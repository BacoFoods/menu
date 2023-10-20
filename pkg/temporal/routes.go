package temporal

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
	routes.GET("/temporal/locales", r.handler.GetLocales)
	routes.GET("/temporal/arqueo", r.handler.GetArqueo)
}
