package order

import (
	"github.com/gin-gonic/gin"
)

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func (r Routes) RegisterRoutes(private, public *gin.RouterGroup) {
	// Order
	private.GET("/order", r.handler.Find)
	private.POST("/order", r.handler.Create)
	public.POST("/order", r.handler.CreatePublic)
	private.GET("/order/:id", r.handler.Get)
	public.GET("/order/:id", r.handler.GetPublic)

	private.PATCH("/order", r.handler.Update)
	private.PATCH("/order/:id/table/:table", r.handler.UpdateTable)
	private.PATCH("/order/:id/seats", r.handler.UpdateSeats)
	private.PATCH("/order/:id/add/products", r.handler.AddProducts)
	public.PATCH("/order/:id/add/products", r.handler.AddProducts)
	private.PATCH("/order/:id/remove/product", r.handler.RemoveProduct)
	private.PATCH("/order/:id/update/product", r.handler.UpdateProduct)
	private.PATCH("/order/:id/update/status", r.handler.UpdateStatus)
	private.PATCH("/order/:id/update/comments", r.handler.UpdateComments)
	private.PATCH("/order/:id/update/client-name", r.handler.UpdateClientName)
	private.POST("/order/:id/release-table", r.handler.ReleaseTable)

	// Order Item
	private.PATCH("/order-item/:id/add/modifiers", r.handler.AddModifiers)
	private.PATCH("/order-item/:id/remove/modifier", r.handler.RemoveModifiers)
	private.PATCH("/order-item/:id/update", r.handler.OrderItemUpdate)

	// Order Types
	private.GET("order-type", r.handler.FindOrderType)
	private.GET("order-type/:id", r.handler.GetOrderType)
	private.POST("order-type", r.handler.CreateOrderType)
	private.PATCH("order-type/:id", r.handler.UpdateOrderType)
	private.DELETE("order-type/:id", r.handler.DeleteOrderType)

	// Invoice
	private.POST("/order/:id/invoice", r.handler.CreateInvoice)
	private.POST("/order/:id/invoice/calculate", r.handler.CalculateInvoice)
	public.GET("/order/:id/invoice/calculate", r.handler.PublicCalculateInvoice)
}
