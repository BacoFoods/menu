package order

import (
	"time"

	"github.com/BacoFoods/menu/internal/telemetry"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) Routes {
	return Routes{handler}
}

func telemetryMiddleware(c *gin.Context) {
	var dt *telemetry.TelemetryPoint
	start := time.Now()
	if c.Request != nil && c.Request.URL != nil {
		dt = telemetry.StartResponse(c.Request.Method, c.Request.URL.Path, time.Now())
	}

	c.Next()

	if c.Writer != nil && dt != nil {
		go func() { dt.Done(c.Writer.Status(), time.Now().Sub(start).Microseconds()) }()
	}
}

func (r Routes) RegisterRoutes(private, public *shared.CustomRoutes) {
	// Order
	private.GET("/order", r.handler.Find, telemetryMiddleware)
	private.POST("/order", r.handler.Create, telemetryMiddleware)
	public.POST("/order", r.handler.CreatePublic, telemetryMiddleware)
	private.GET("/order/:id", r.handler.Get, telemetryMiddleware)
	public.GET("/order/:id", r.handler.GetPublic, telemetryMiddleware)

	private.PATCH("/order", r.handler.Update)
	private.PATCH("/order/:id/table/:table", r.handler.UpdateTable)
	private.PATCH("/order/:id/seats", r.handler.UpdateSeats)
	private.PATCH("/order/:id/add/products", r.handler.AddProducts)
	public.PATCH("/order/:id/add/products", r.handler.AddProducts)
	private.PATCH("/order/:id/remove/product", r.handler.RemoveProduct)
	private.PATCH("/order/:id/update/product", r.handler.UpdateProduct)
	private.PATCH("/order/:id/update/comments", r.handler.UpdateComments)
	private.PATCH("/order/:id/update/client-name", r.handler.UpdateClientName)
	private.PATCH("/order/:id/update/status", r.handler.UpdateStatus)

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
	private.POST("/order/:id/invoice", r.handler.CreateInvoice, telemetryMiddleware)
	private.POST("/order/:id/invoice/calculate", r.handler.CalculateInvoice)
	public.GET("/order/:id/invoice/calculate", r.handler.PublicCalculateInvoice)
	public.POST("/order/:id/checkout", r.handler.PublicCheckout)
	private.POST("/invoice/:id/close", r.handler.CloseInvoice)
}
