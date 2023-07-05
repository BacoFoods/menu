package brand

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r *Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/brand", r.handler.Find)
	router.GET("/brand/:id", r.handler.Get)
	router.POST("/brand", r.handler.Create)
	router.PATCH("/brand/:id", r.handler.Update)
	router.DELETE("/brand/:id", r.handler.Delete)
}
