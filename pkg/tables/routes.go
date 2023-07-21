package tables

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/tables", r.handler.Find)
	router.GET("/tables/:id", r.handler.Get)
	router.POST("/tables", r.handler.Create)
	router.PATCH("/tables/:id", r.handler.Update)
	router.DELETE("/tables/:id", r.handler.Delete)
}
