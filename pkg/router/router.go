package router

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/availability"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/country"
	"github.com/BacoFoods/menu/pkg/currency"
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/overriders"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/surcharge"
	"github.com/BacoFoods/menu/pkg/swagger"
	"github.com/BacoFoods/menu/pkg/tables"
	"github.com/BacoFoods/menu/pkg/taxes"
	"github.com/BacoFoods/menu/pkg/zones"
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
	routes.Availability.RegisterRoutes(private)
	routes.Brand.RegisterRoutes(private)
	routes.Category.RegisterRoutes(private)
	routes.Channel.RegisterRoutes(private)
	routes.Country.RegisterRoutes(private)
	routes.Currency.RegisterRoutes(private)
	routes.Discount.RegisterRoutes(private)
	routes.Menu.RegisterRoutes(private)
	routes.Order.RegisterRoutes(private)
	routes.Overriders.RegisterRoutes(private)
	routes.Product.RegisterRoutes(private)
	routes.Store.RegisterRoutes(private)
	routes.Surcharge.RegisterRoutes(private)
	routes.Taxes.RegisterRoutes(private)
	routes.Table.RegisterRoutes(private)
	routes.Zone.RegisterRoutes(private)

	// Register public routes
	public := router.Group(fmt.Sprintf("%s/public", path))
	routes.Swagger.Register(public)

	return router
}

// RoutesGroup for unify all routes
type RoutesGroup struct {
	Availability availability.Routes
	Brand        brand.Routes
	Category     category.Routes
	Channel      channel.Routes
	Country      country.Routes
	Currency     currency.Routes
	Discount     discount.Routes
	HealthCheck  healthcheck.Routes
	Menu         menu.Routes
	Overriders   overriders.Routes
	Order        order.Routes
	Product      product.Routes
	Store        store.Routes
	Surcharge    surcharge.Routes
	Swagger      swagger.Routes
	Table        tables.Routes
	Taxes        taxes.Routes
	Zone         zones.Routes
}
