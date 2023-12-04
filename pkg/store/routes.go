package store

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *shared.CustomRoutes) {
	private.GET("/store", r.handler.Find)
	private.GET("/store/:id", r.handler.Get)
	private.POST("/store", r.handler.Create)
	private.PATCH("/store/:id", r.handler.Update)
	private.DELETE("/store/:id", r.handler.Delete)

	private.PATCH("/store/:id/enable", r.handler.Enable)
	private.PATCH("/store/:id/channel/:channelID", r.handler.AddChannel)

	private.GET("/store/:id/zone", r.handler.FindZonesByStore)
	private.GET("/store/:id/zone/:zoneID", r.handler.GetZoneByStore)
}
