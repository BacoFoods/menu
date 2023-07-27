package menu

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup) {
	private.GET("/menu", r.handler.Find)
	private.GET("/menu/:id", r.handler.Get)
	private.GET("/menu/place/:place/:place-id/list", r.handler.ListByPlace)
	private.GET("/menu/place/:place/:place-id/menu-id/:menu-id", r.handler.GetByPlace)
	private.GET("/menu/:id/store/:storeID/channels", r.handler.FindChannels)
	private.POST("/menu", r.handler.Create)
	private.PUT("/menu/:id/place/:place/availability", r.handler.UpdateAvailability)
	private.PATCH("/menu/:id", r.handler.Update)
	private.PATCH("/menu/:id/category/:categoryID/add", r.handler.AddCategory)
	private.PATCH("/menu/:id/category/:categoryID/remove", r.handler.RemoveCategory)
	private.DELETE("/menu/:id", r.handler.Delete)
}
