package main

import (
	"fmt"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/country"
	"github.com/BacoFoods/menu/pkg/currency"
	"github.com/BacoFoods/menu/pkg/database"
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/router"
	"github.com/BacoFoods/menu/pkg/swagger"
	"github.com/BacoFoods/menu/pkg/taxes"
	"github.com/sirupsen/logrus"
)

func main() {
	// Database
	gormFramework := database.MustNewGormFramework()
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
	)

	// Healthcheck
	healthcheckHandler := healthcheck.NewHandler()
	healthcheckRoutes := healthcheck.NewRoutes(healthcheckHandler)

	// Swagger
	swaggerRoutes := swagger.NewRoutes()

	// Menu
	menuRepository := menu.NewGormRepository(gormDB)
	menuService := menu.NewService(menuRepository)
	menuHandler := menu.NewHandler(menuService)
	menuRoutes := menu.NewRoutes(menuHandler)

	// Category
	categoryRepository := category.NewGormRepository(gormDB)
	categoryService := category.NewService(categoryRepository)
	categoryHandler := category.NewHandler(categoryService)
	categoryRoutes := category.NewRoutes(categoryHandler)

	// Product
	productRepository := product.NewGormRepository(gormDB)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService)
	productRoutes := product.NewRoutes(productHandler)

	// Taxes
	taxesRepository := taxes.NewGormRepository(gormDB)
	taxesService := taxes.NewService(taxesRepository)
	taxesHandler := taxes.NewHandler(taxesService)
	taxesRoutes := taxes.NewRoutes(taxesHandler)

	// Country
	countryRepository := country.NewGormRepository(gormDB)
	countryService := country.NewService(countryRepository)
	countryHandler := country.NewHandler(countryService)
	countryRoutes := country.NewRoutes(countryHandler)

	// Currency
	currencyRepository := currency.NewGormRepository(gormDB)
	currencyService := currency.NewService(currencyRepository)
	currencyHandler := currency.NewHandler(currencyService)
	currencyRoutes := currency.NewRoutes(currencyHandler)

	// Routes
	routes := &router.RoutesGroup{
		HealthCheck: healthcheckRoutes,
		Swagger:     swaggerRoutes,
		Menu:        menuRoutes,
		Category:    categoryRoutes,
		Product:     productRoutes,
		Taxes:       taxesRoutes,
		Country:     countryRoutes,
		Currency:    currencyRoutes,
	}

	// Run server
	r := router.NewRouter(routes)
	logrus.Fatal(r.Run(fmt.Sprintf(":%v", internal.Config.AppPort)))
}
