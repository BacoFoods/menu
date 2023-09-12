package router

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/account"
	"github.com/BacoFoods/menu/pkg/availability"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/country"
	"github.com/BacoFoods/menu/pkg/currency"
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/BacoFoods/menu/pkg/status"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/surcharge"
	"github.com/BacoFoods/menu/pkg/swagger"
	"github.com/BacoFoods/menu/pkg/tables"
	"github.com/BacoFoods/menu/pkg/taxes"
	"github.com/BacoFoods/menu/pkg/zones"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

// Router interface for router implementation
type Router interface {
	Run(addr ...string) error
}

// NewRouter create a new router instance with all routes using gin
func NewRouter(routes *RoutesGroup) Router {
	path := "api/menu/v1"
	fmt.Println(internal.Config.GoogleConfig)
	decodeBytes, err := base64.StdEncoding.DecodeString(internal.Config.GoogleConfig)
	if err != nil {
		shared.LogError("error decoding jwt key", "pkg/router/router.go", "NewRouter", err, nil)
		panic(err)
	}

	validator, err := idtoken.NewValidator(context.TODO(), option.WithCredentialsJSON(decodeBytes))
	if err != nil {
		shared.LogError("error initializing validator", "pkg/router/router.go", "NewRouter", err, nil)
	}

	// Setting middlewares
	router := gin.Default()
	router.Use(CORSMiddleware())

	// Register health check route
	healthCheck := router.Group(path)

	routes.HealthCheck.Register(healthCheck)

	// Register private routes
	private := router.Group(path)
	private.Use(AuthMiddleware(validator))

	routes.Availability.RegisterRoutes(private)
	routes.Brand.RegisterRoutes(private)
	routes.Category.RegisterRoutes(private)
	routes.Channel.RegisterRoutes(private)
	routes.Country.RegisterRoutes(private)
	routes.Currency.RegisterRoutes(private)
	routes.Discount.RegisterRoutes(private)
	routes.Menu.RegisterRoutes(private)
	routes.Order.RegisterRoutes(private)
	routes.Status.RegisterRoutes(private)
	routes.Product.RegisterRoutes(private)
	routes.Store.RegisterRoutes(private)
	routes.Surcharge.RegisterRoutes(private)
	routes.Taxes.RegisterRoutes(private)
	routes.Table.RegisterRoutes(private)
	routes.Zone.RegisterRoutes(private)
	routes.Invoice.RegisterRoutes(private)

	// Register public routes
	public := router.Group(fmt.Sprintf("%s/public", path))

	routes.Account.RegisterRoutes(private, public)
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
	Order        order.Routes
	Status       status.Routes
	Product      product.Routes
	Store        store.Routes
	Surcharge    surcharge.Routes
	Swagger      swagger.Routes
	Table        tables.Routes
	Taxes        taxes.Routes
	Zone         zones.Routes
	Invoice      invoice.Routes
	Account      account.Routes
}
