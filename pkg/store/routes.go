package store

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/store", r.handler.Find)
	router.GET("/store/:id", r.handler.Get)
	router.POST("/store", r.handler.Create)
	router.PATCH("/store/:id", r.handler.Update)
	router.DELETE("/store/:id", r.handler.Delete)

	router.PATCH("/store/:id/channel/:channelID", r.handler.AddChannel)

	router.GET("/store/:id/zone", r.handler.FindZonesByStore)
	router.GET("/store/:id/zone/:zoneID", r.handler.GetZoneByStore)

}
