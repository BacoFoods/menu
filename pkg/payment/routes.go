package payment

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/payment/:id", r.handler.Get)
	router.GET("/payment", r.handler.Find)
	router.POST("/payment", r.handler.Create)
	router.PATCH("/payment", r.handler.Update)
	router.DELETE("/payment/:id", r.handler.Delete)
}
