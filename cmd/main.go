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
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/overriders"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/router"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/swagger"
	"github.com/BacoFoods/menu/pkg/taxes"
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
		&category.CategoriesProducts{},
		&product.Product{},
		&taxes.Tax{},
		&country.Country{},
		&currency.Currency{},
		&brand.Brand{},
		&store.Store{},
		&channel.Channel{},
		&overriders.Overriders{},
		&availability.Availability{},
	)

	// Healthcheck
	healthcheckHandler := healthcheck.NewHandler()
	healthcheckRoutes := healthcheck.NewRoutes(healthcheckHandler)

	// Swagger
	swaggerRoutes := swagger.NewRoutes()

	// Overriders
	overridersRepository := overriders.NewDBRepository(gormDB)
	overridersService := overriders.NewService(overridersRepository)
	overridersHandler := overriders.NewHandler(overridersService)
	overridersRoutes := overriders.NewRoutes(overridersHandler)

	// Availability
	availabilityRepository := availability.NewDBRepository(gormDB)
	availabilityService := availability.NewService(availabilityRepository)
	availabilityHandler := availability.NewHandler(availabilityService)
	availabilityRoutes := availability.NewRoutes(availabilityHandler)

	// Menu
	menuRepository := menu.NewDBRepository(gormDB)
	menuService := menu.NewService(menuRepository, overridersRepository, availabilityRepository)
	menuHandler := menu.NewHandler(menuService)
	menuRoutes := menu.NewRoutes(menuHandler)

	// Category
	categoryRepository := category.NewDBRepository(gormDB)
	categoryService := category.NewService(categoryRepository)
	categoryHandler := category.NewHandler(categoryService)
	categoryRoutes := category.NewRoutes(categoryHandler)

	// Product
	productRepository := product.NewDBRepository(gormDB)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService)
	productRoutes := product.NewRoutes(productHandler)

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

	// Store
	storeRepository := store.NewDBRepository(gormDB)
	storeService := store.NewService(storeRepository)
	storeHandler := store.NewHandler(storeService)
	storeRoutes := store.NewRoutes(storeHandler)

	// Channel
	channelRepository := channel.NewDBRepository(gormDB)
	channelService := channel.NewService(channelRepository)
	channelHandler := channel.NewHandler(channelService)
	channelRoutes := channel.NewRoutes(channelHandler)

	// Routes
	routes := &router.RoutesGroup{
		HealthCheck:  healthcheckRoutes,
		Swagger:      swaggerRoutes,
		Menu:         menuRoutes,
		Category:     categoryRoutes,
		Product:      productRoutes,
		Overriders:   overridersRoutes,
		Taxes:        taxesRoutes,
		Country:      countryRoutes,
		Currency:     currencyRoutes,
		Brand:        brandRoutes,
		Store:        storeRoutes,
		Channel:      channelRoutes,
		Availability: availabilityRoutes,
	}

	// Run server
	r := router.NewRouter(routes)
	logrus.Fatal(r.Run(fmt.Sprintf(":%v", internal.Config.AppPort)))
}
