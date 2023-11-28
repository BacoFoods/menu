package account

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private, public *shared.CustomRoutes) {
	public.POST("/account", r.handler.Create)
	public.POST("/account/login", r.handler.Login)
	public.POST("/account/login/pin", r.handler.LoginPin)
	private.POST("/account", r.handler.CreatePinUser)
	private.DELETE("/account/:id", r.handler.Delete)
	private.GET("/account", r.handler.Find)
	private.PATCH("/account", r.handler.Update)
}
