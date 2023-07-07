package router

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/availability"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/country"
	"github.com/BacoFoods/menu/pkg/currency"
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/overriders"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/swagger"
	"github.com/BacoFoods/menu/pkg/taxes"
	"github.com/gin-gonic/gin"
)

// Router interface for router implementation
type Router interface {
	Run(addr ...string) error
}

// NewRouter create a new router instance with all routes using gin
func NewRouter(routes *RoutesGroup) Router {
	path := "api/menu/v1"

	// Setting middlewares
	router := gin.Default()
	router.Use(CORSMiddleware(), Authentication())

	// Register health check route
	healthCheck := router.Group(path)
	routes.HealthCheck.Register(healthCheck)

	// Register private routes
	private := router.Group(path)
	routes.Menu.RegisterRoutes(private)
	routes.Category.RegisterRoutes(private)
	routes.Product.RegisterRoutes(private)
	routes.Taxes.RegisterRoutes(private)
	routes.Country.RegisterRoutes(private)
	routes.Currency.RegisterRoutes(private)
	routes.Brand.RegisterRoutes(private)
	routes.Store.RegisterRoutes(private)
	routes.Channel.RegisterRoutes(private)
	routes.Availability.RegisterRoutes(private)

	// Register public routes
	public := router.Group(fmt.Sprintf("%s/public", path))
	routes.Swagger.Register(public)

	return router
}

// RoutesGroup for unify all routes
type RoutesGroup struct {
	HealthCheck  healthcheck.Routes
	Swagger      swagger.Routes
	Menu         menu.Routes
	Category     category.Routes
	Product      product.Routes
	Overriders   overriders.Routes
	Taxes        taxes.Routes
	Country      country.Routes
	Currency     currency.Routes
	Brand        brand.Routes
	Store        store.Routes
	Channel      channel.Routes
	Availability availability.Routes
}
