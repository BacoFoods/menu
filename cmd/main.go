package main

import (
	"fmt"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/database"
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/router"
	"github.com/BacoFoods/menu/pkg/swagger"
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

	// Routes
	routes := &router.RoutesGroup{
		HealthCheck: healthcheckRoutes,
		Swagger:     swaggerRoutes,
		Menu:        menuRoutes,
	}

	// Run server
	r := router.NewRouter(routes)
	logrus.Fatal(r.Run(fmt.Sprintf(":%v", internal.Config.AppPort)))
}
