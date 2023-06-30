package menu

import "github.com/gin-gonic/gin"

// Routes is the struct that contains all the routes for the menu package
type Routes struct {
	handler *Handler
}

// NewRoutes creates a new instance of the Routes struct
func NewRoutes(handler *Handler) Routes {
	return Routes{
		handler: handler,
	}
}

// Register registers all the routes for the menu package
func (r *Routes) Register(private *gin.RouterGroup) {
	private.POST("/menu", r.handler.Create)
	private.GET("/menu", r.handler.Find)
	private.GET("/menu/:id", r.handler.Get)
	private.PATCH("/menu", r.handler.Update)
	private.DELETE("/menu/:id", r.handler.Delete)
}
