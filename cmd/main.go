package main

import (
	"fmt"

	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/availability"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/country"
	"github.com/BacoFoods/menu/pkg/currency"
	"github.com/BacoFoods/menu/pkg/database"
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/router"
	"github.com/BacoFoods/menu/pkg/status"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/surcharge"
	"github.com/BacoFoods/menu/pkg/swagger"
	"github.com/BacoFoods/menu/pkg/tables"
	"github.com/BacoFoods/menu/pkg/taxes"
	"github.com/BacoFoods/menu/pkg/zones"
	"github.com/sirupsen/logrus"
)

func main() {
	// Database
	gormFramework := database.MustNewGormFramework("")
	gormDB := gormFramework.GetDBClient()

	// DB Migrations
	gormFramework.MustMakeMigrations(
		&menu.Menu{},
		&menu.MenusCategories{},
		&category.Category{},
		&discount.Discount{},
		&surcharge.Surcharge{},
		&product.Product{},
		&product.Modifier{},
		&product.Overrider{},
		&taxes.Tax{},
		&country.Country{},
		&currency.Currency{},
		&brand.Brand{},
		&store.Store{},
		&zones.Zone{},
		&tables.Table{},
		&channel.Channel{},
		&availability.Availability{},
		&order.Order{},
		&order.OrderItem{},
		&order.OrderModifier{},
		&order.OrderType{},
		&status.Status{},
		&invoice.Invoice{},
		&invoice.Item{},
		&invoice.Discount{},
		&invoice.Surcharge{},
	)

	// Healthcheck
	healthcheckHandler := healthcheck.NewHandler()
	healthcheckRoutes := healthcheck.NewRoutes(healthcheckHandler)

	// Swagger
	swaggerRoutes := swagger.NewRoutes()

	// Channel
	channelRepository := channel.NewDBRepository(gormDB)
	channelService := channel.NewService(channelRepository)
	channelHandler := channel.NewHandler(channelService)
	channelRoutes := channel.NewRoutes(channelHandler)

	// Store
	storeRepository := store.NewDBRepository(gormDB)
	storeService := store.NewService(storeRepository, channelRepository)
	storeHandler := store.NewHandler(storeService)
	storeRoutes := store.NewRoutes(storeHandler)

	// Zone
	zoneRepository := zones.NewDBRepository(gormDB)
	zoneService := zones.NewService(zoneRepository)
	zoneHandler := zones.NewHandler(zoneService)
	zoneRoutes := zones.NewRoutes(zoneHandler)

	// Tables
	tablesRepository := tables.NewDBRepository(gormDB)
	tablesService := tables.NewService(tablesRepository)
	tablesHandler := tables.NewHandler(tablesService)
	tablesRoutes := tables.NewRoutes(tablesHandler)

	// Availability
	availabilityRepository := availability.NewDBRepository(gormDB)
	availabilityService := availability.NewService(availabilityRepository, storeRepository, channelRepository)
	availabilityHandler := availability.NewHandler(availabilityService)
	availabilityRoutes := availability.NewRoutes(availabilityHandler)

	// Product
	productRepository := product.NewDBRepository(gormDB)
	productService := product.NewService(productRepository, channelRepository)
	productHandler := product.NewHandler(productService)
	productRoutes := product.NewRoutes(productHandler)

	// Surcharge
	surchargeRepository := surcharge.NewDBRepository(gormDB)
	surchargeService := surcharge.NewService(surchargeRepository)
	surchargeHandler := surcharge.NewHandler(surchargeService)
	surchargeRoutes := surcharge.NewRoutes(surchargeHandler)

	// Discount
	discountRepository := discount.NewDBRepository(gormDB)
	discountService := discount.NewService(discountRepository)
	discountHandler := discount.NewHandler(discountService)
	discountRoutes := discount.NewRoutes(discountHandler)

	// Category
	categoryRepository := category.NewDBRepository(gormDB)
	categoryService := category.NewService(categoryRepository, productRepository)
	categoryHandler := category.NewHandler(categoryService)
	categoryRoutes := category.NewRoutes(categoryHandler)

	// Menu
	menuRepository := menu.NewDBRepository(gormDB)
	menuService := menu.NewService(menuRepository, productRepository, availabilityRepository, storeRepository, categoryRepository)
	menuHandler := menu.NewHandler(menuService)
	menuRoutes := menu.NewRoutes(menuHandler)

	// Taxes
	taxesRepository := taxes.NewDBRepository(gormDB)
	taxesService := taxes.NewService(taxesRepository)
	taxesHandler := taxes.NewHandler(taxesService)
	taxesRoutes := taxes.NewRoutes(taxesHandler)

	// Country
	countryRepository := country.NewDBRepository(gormDB)
	countryService := country.NewService(countryRepository)
	countryHandler := country.NewHandler(countryService)
	countryRoutes := country.NewRoutes(countryHandler)

	// Currency
	currencyRepository := currency.NewDBRepository(gormDB)
	currencyService := currency.NewService(currencyRepository)
	currencyHandler := currency.NewHandler(currencyService)
	currencyRoutes := currency.NewRoutes(currencyHandler)

	// Brand
	brandRepository := brand.NewDBRepository(gormDB)
	brandService := brand.NewService(brandRepository)
	brandHandler := brand.NewHandler(brandService)
	brandRoutes := brand.NewRoutes(brandHandler)

	// Status
	statusRepository := status.NewDBRepository(gormDB)
	statusService := status.NewService(statusRepository)
	statusHandler := status.NewHandler(statusService)
	statusRoutes := status.NewRoutes(statusHandler)

	// Invoice
	invoiceRepository := invoice.NewDBRepository(gormDB)
	invoiceService := invoice.NewService(invoiceRepository)
	invoiceHandler := invoice.NewHandler(invoiceService)
	invoiceRoutes := invoice.NewRoutes(invoiceHandler)

	// Order
	orderRepository := order.NewDBRepository(gormDB)
	orderService := order.NewService(orderRepository, tablesRepository, productRepository, invoiceRepository)
	orderHandler := order.NewHandler(orderService)
	orderRoutes := order.NewRoutes(orderHandler)

	// Routes
	routes := &router.RoutesGroup{
		HealthCheck:  healthcheckRoutes,
		Swagger:      swaggerRoutes,
		Menu:         menuRoutes,
		Category:     categoryRoutes,
		Product:      productRoutes,
		Surcharge:    surchargeRoutes,
		Discount:     discountRoutes,
		Taxes:        taxesRoutes,
		Country:      countryRoutes,
		Currency:     currencyRoutes,
		Brand:        brandRoutes,
		Store:        storeRoutes,
		Zone:         zoneRoutes,
		Table:        tablesRoutes,
		Channel:      channelRoutes,
		Availability: availabilityRoutes,
		Order:        orderRoutes,
		Status:       statusRoutes,
		Invoice: 	  *invoiceRoutes,
	}

	// Run server
	r := router.NewRouter(routes)
	logrus.Fatal(r.Run(fmt.Sprintf(":%v", internal.Config.AppPort)))
}
