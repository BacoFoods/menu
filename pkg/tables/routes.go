package tables

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private *gin.RouterGroup, public *gin.RouterGroup) {
	private.GET("/tables", r.handler.Find)
	private.GET("/tables/:id", r.handler.Get)
	private.POST("/tables/:id/generate", r.handler.GenerateQR)
	private.POST("/tables", r.handler.Create)
	private.PATCH("/tables/:id", r.handler.Update)
	private.DELETE("/tables/:id", r.handler.Delete)

	public.GET("/tables/scan/:qrId", r.handler.ScanQR)
}
