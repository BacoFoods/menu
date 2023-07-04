package spot

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/spot", r.handler.Find)
	router.GET("/spot/:id", r.handler.Get)
	router.POST("/spot", r.handler.Create)
	router.PATCH("/spot/:id", r.handler.Update)
	router.DELETE("/spot/:id", r.handler.Delete)
}
