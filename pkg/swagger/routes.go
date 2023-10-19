package swagger

import (
	"github.com/BacoFoods/menu/pkg/account"
	"github.com/BacoFoods/menu/pkg/availability"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/client"
	"github.com/BacoFoods/menu/pkg/course"
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/payment"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/shift"
	"github.com/BacoFoods/menu/pkg/status"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/surcharge"
	_ "github.com/BacoFoods/menu/pkg/swagger/docs" //documentation
	"github.com/BacoFoods/menu/pkg/zones"
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
// @contact.email andersonrodriguezce@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/menu/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
//
//go:generate go run github.com/swaggo/swag/cmd/swag@v1.8.4 init --parseDependency=true -o docs -g routes.go
func (r *Routes) Register(group *gin.RouterGroup) {
	// Next line is for swagger documentation
	_ = menu.Handler{}
	_ = category.Handler{}
	_ = product.Handler{}
	_ = discount.Handler{}
	_ = surcharge.Handler{}
	_ = brand.Handler{}
	_ = store.Handler{}
	_ = zones.Handler{}
	_ = channel.Handler{}
	_ = availability.Handler{}
	_ = order.Handler{}
	_ = status.Handler{}
	_ = invoice.Handler{}
	_ = account.Handler{}
	_ = course.Handler{}
	_ = client.Handler{}
	_ = payment.Handler{}
	_ = shift.Handler{}
	group.GET("/swagger/*any", swaggerGin.WrapHandler(swaggerFiles.Handler))
}
