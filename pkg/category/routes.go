package category

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/category", r.handler.Find)
	router.GET("/category/:id", r.handler.Get)
	router.POST("/category", r.handler.Create)
	router.PATCH("/category/:id", r.handler.Update)
	router.DELETE("/category/:id", r.handler.Delete)
}