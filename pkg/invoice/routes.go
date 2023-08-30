package invoice

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/invoice/:id", r.handler.Get)
	router.PATCH("/invoice/:id", r.handler.Update)
}
