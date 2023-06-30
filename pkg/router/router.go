package router

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/swagger"
	"github.com/gin-gonic/gin"
)

// Router interface for router implementation
type Router interface {
	Run(addr ...string) error
}

// NewRouter create a new router instance with all routes using gin
func NewRouter(routes *RoutesGroup) Router {
	path := "api/menu/v1"

	router := gin.Default()
	router.Use(CORSMiddleware(), Authentication())

	healthCheck := router.Group(path)
	routes.HealthCheck.Register(healthCheck)

	private := router.Group(path)
	routes.Menu.Register(private)

	public := router.Group(fmt.Sprintf("%s/public", path))
	routes.Swagger.Register(public)

	return router
}

// RoutesGroup for unify all routes
type RoutesGroup struct {
	HealthCheck healthcheck.Routes
	Swagger     swagger.Routes
	Menu        menu.Routes
}
