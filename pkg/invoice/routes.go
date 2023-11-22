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
	router.GET("/invoice/:id/print", r.handler.Print)
	router.POST("/invoice/:id", r.handler.UpdateTip)
	router.POST("/invoice/:id/client/:clientID/add", r.handler.AddClient)
	router.POST("/invoice/:id/client/:clientID/remove", r.handler.RemoveClient)
	router.POST("/invoice/:id/split", r.handler.Split)

	router.GET("/discount-applied", r.handler.FindDiscountApplied)
	router.DELETE("/discount-applied/:id", r.handler.RemoveDiscountApplied)
}
