package swagger

import (
	_ "github.com/BacoFoods/menu/pkg/swagger/docs" //documentation

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	swaggerGin "github.com/swaggo/gin-swagger"
)

// Routes for swagger routes
type Routes struct{}

// NewRoutes create a new swagger routes instance
func NewRoutes() Routes {
	return Routes{}
}

// Register to register the routes
// Next line is for swagger documentation
// @title MENU-MS API Rest
// @version 1.0
// @description Microservices for MENU.
// @termsOfService https://www.bacu.co/
// @contact.name Anderson Rodriguez
// @contact.url https://www.bacu.co/quienes-somos
// @contact.email anderson.rodriguez@bacu.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/menu/v1
//
//go:generate go run github.com/swaggo/swag/cmd/swag@v1.8.4 init --parseDependency=true -o docs -g routes.go
func (r *Routes) Register(group *gin.RouterGroup) {
	// Next line is for swagger documentation
	// _ = menu.Handler{}

	group.GET("/swagger/*any", swaggerGin.WrapHandler(swaggerFiles.Handler))
}