package product

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/product", r.handler.Find)
	router.GET("/product/:id", r.handler.Get)
	router.POST("/product", r.handler.Create)
	router.PATCH("/product/:id", r.handler.Update)
	router.DELETE("/product/:id", r.handler.Delete)
}
