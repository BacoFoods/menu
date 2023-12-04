package router

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/BacoFoods/menu/pkg/app"
	"github.com/BacoFoods/menu/pkg/siesa"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"

	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/account"
	"github.com/BacoFoods/menu/pkg/assets"
	"github.com/BacoFoods/menu/pkg/availability"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/cashaudit"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/client"
	"github.com/BacoFoods/menu/pkg/connector"
	"github.com/BacoFoods/menu/pkg/country"
	"github.com/BacoFoods/menu/pkg/course"
	"github.com/BacoFoods/menu/pkg/currency"
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/facturacion"
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/payment"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/scheduler"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/BacoFoods/menu/pkg/shift"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/surcharge"
	"github.com/BacoFoods/menu/pkg/swagger"
	"github.com/BacoFoods/menu/pkg/tables"
	"github.com/BacoFoods/menu/pkg/taxes"
	"github.com/BacoFoods/menu/pkg/temporal"
	"github.com/gin-gonic/gin"
)

// Router interface for router implementation
type Router interface {
	Run(addr ...string) error
}

// NewRouter create a new router instance with all routes using gin
func NewRouter(routes *RoutesGroup) Router {
	path := "api/menu/v1"

	// Setting Google JWT validator
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
	routes.Product.RegisterRoutes(private)
	routes.Store.RegisterRoutes(private)
	routes.Surcharge.RegisterRoutes(private)
	routes.Taxes.RegisterRoutes(private)
	routes.Invoice.RegisterRoutes(private)
	routes.Course.RegisterRoutes(private)
	routes.Client.RegisterRoutes(private)
	routes.Payment.RegisterRoutes(private)
	routes.Cashier.RegisterRoutes(private)
	routes.Temporal.RegisterRoutes(private)
	routes.CashAudit.RegisterRoutes(private)
	routes.Assets.RegisterRoutes(private)
	routes.Facturacion.RegisterRoutes(private)
	routes.Schedule.RegisterRoutes(private)
	routes.Equivalence.RegisterRoutes(private)
	routes.Siesa.RegisterRoutes(private)
	routes.App.RegisterRoutes(private)

	// Register public routes
	public := router.Group(fmt.Sprintf("%s/public", path))

	routes.Table.RegisterRoutes(private, public)
	routes.Menu.RegisterRoutes(private, public)
	routes.Account.RegisterRoutes(private, public)
	routes.Order.RegisterRoutes(private, public)

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
	Product      product.Routes
	Store        store.Routes
	Surcharge    surcharge.Routes
	Swagger      swagger.Routes
	Table        tables.Routes
	Taxes        taxes.Routes
	Invoice      invoice.Routes
	Account      account.Routes
	Course       course.Routes
	Client       client.Routes
	Payment      payment.Routes
	Cashier      shift.Routes
	Temporal     temporal.Routes
	CashAudit    cashaudit.Routes
	Assets       assets.Routes
	Facturacion  facturacion.Routes
	Schedule     scheduler.Routes
	Equivalence  connector.Routes
	Siesa        siesa.Routes
	App          app.Routes
}
