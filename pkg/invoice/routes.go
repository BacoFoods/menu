package invoice

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler: handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/invoice/:id", r.handler.Get)
	router.GET("/invoice", r.handler.Find)
	router.POST("/invoice/:id", r.handler.UpdateTip)
	router.POST("/invoice/:id/client/:clientID/add", r.handler.AddClient)
	router.POST("/invoice/:id/client/:clientID/remove", r.handler.RemoveClient)
	router.POST("/invoice/:id/separate", r.handler.Separate)
}
