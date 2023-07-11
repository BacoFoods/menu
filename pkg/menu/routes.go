package menu

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/menu", r.handler.Find)
	router.GET("/menu/:id", r.handler.Get)
	router.GET("/menu/place/:place/:place-id/list", r.handler.ListByPlace)
	router.GET("/menu/place/:place/:place-id/menu-id/:menu-id", r.handler.GetByPlace)
	router.GET("/menu/:id/store/:storeID/channels", r.handler.FindChannels)
	router.POST("/menu", r.handler.Create)
	router.PUT("/menu/:id/place/:place/availability", r.handler.UpdateAvailability)
	router.PATCH("/menu/:id/category/:categoryID", r.handler.AddCategory)
	router.PATCH("/menu/:id", r.handler.Update)
	router.DELETE("/menu/:id", r.handler.Delete)
}
