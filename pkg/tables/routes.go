package tables

import "github.com/BacoFoods/menu/pkg/shared"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private, public *shared.CustomRoutes) {
	private.GET("/tables", r.handler.Find)
	private.GET("/tables/:id", r.handler.Get)
	private.POST("/tables/:id/generate", r.handler.GenerateQR)
	private.POST("/tables", r.handler.Create)
	private.PATCH("/tables/:id", r.handler.Update)
	private.DELETE("/tables/:id", r.handler.Delete)
	private.POST("/tables/:id/release", r.handler.Release)

	private.GET("/zone", r.handler.FindZones)
	private.GET("/zone/:id", r.handler.GetZone)
	private.POST("/zone", r.handler.CreateZone)
	private.PATCH("/zone/:id", r.handler.UpdateZone)
	private.DELETE("/zone/:id", r.handler.DeleteZone)
	private.PATCH("/zone/:id/tables/add", r.handler.AddTables)
	private.PATCH("/zone/:id/tables/remove", r.handler.RemoveTables)
	private.PATCH("/zone/:id/enable", r.handler.Enable)

	public.GET("/tables/scan/:qrId", r.handler.ScanQR)
}
