package zones

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r *Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/zone", r.handler.Find)
	router.GET("/zone/:id", r.handler.Get)
	router.POST("/zone", r.handler.Create)
	router.PATCH("/zone/:id", r.handler.Update)
	router.DELETE("/zone/:id", r.handler.Delete)

	router.PATCH("/zone/:id/tables/add", r.handler.AddTables)
	router.PATCH("/zone/:id/tables/remove", r.handler.RemoveTables)
}
