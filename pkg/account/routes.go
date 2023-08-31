package account

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private, public *gin.RouterGroup) {
	public.POST("/account", r.handler.Create)
	private.POST("/account/login", r.handler.Login)
	private.DELETE("/account/:id", r.handler.Delete)
	private.GET("/account", r.handler.Find)
}
