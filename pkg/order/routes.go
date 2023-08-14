package order

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(router *gin.RouterGroup) {
	// Order
	router.POST("/order", r.handler.Create)
	router.PATCH("/order/:id/table/:table", r.handler.UpdateTable)
	router.PATCH("/order/:id/seats", r.handler.UpdateSeats)
	router.GET("/order/:id", r.handler.Get)
	router.GET("/order", r.handler.Find)

	// Order Types
	router.GET("order-type", r.handler.FindOrderType)
	router.GET("order-type/:id", r.handler.GetOrderType)
	router.POST("order-type", r.handler.CreateOrderType)
	router.PATCH("order-type/:id", r.handler.UpdateOrderType)
	router.DELETE("order-type/:id", r.handler.DeleteOrderType)
}
