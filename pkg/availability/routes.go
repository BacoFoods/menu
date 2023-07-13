package availability

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.PUT("/availability/:entity/:entity-id/:place/:place-id", r.handler.EnableEntity)
	router.GET("/availability/entities", r.handler.FindEntities)
	router.GET("/availability/places", r.handler.FindPlaces)

	router.GET("/availability/:entity/:entity-id/:place", r.handler.Find)
	router.GET("/availability/:entity/:entity-id/:place/:place-id", r.handler.Get)
}
