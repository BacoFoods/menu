package healthcheck

import "github.com/gin-gonic/gin"

// Routes for health check routes with a handler instance as dependency
type Routes struct {
	handler Handler
}

// NewRoutes create a new routes instance with a handler instance as dependency
func NewRoutes(handler Handler) Routes {
	return Routes{handler}
}

// Register all routes for health check associating route to handler
func (r *Routes) Register(group *gin.RouterGroup) {
	group.GET("/health-check", r.handler.HealthCheck)
}
