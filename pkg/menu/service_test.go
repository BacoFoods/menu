package menu_test

import (
	"fmt"
	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/country"
	"github.com/BacoFoods/menu/pkg/currency"
	"github.com/BacoFoods/menu/pkg/database"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/overriders"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/taxes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"time"
)

var _ = Describe("Service", func() {
	var (
		db          *gorm.DB
		menuService menu.Service
	)

	BeforeSuite(func() {
		// Database
		customDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			internal.Config.DBConfig.Host,
			internal.Config.DBConfig.User,
			internal.Config.DBConfig.Password,
			internal.Config.DBConfig.Name,
			internal.Config.DBConfig.Port,
		)
		gormFramework := database.MustNewGormFramework(customDSN)
		db = gormFramework.GetDBClient()

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
		)

		// Service Implementation
		menuRepository := menu.NewDBRepository(db)
		overriderRepository := overriders.NewDBRepository(db)
		menuService = menu.NewService(menuRepository, overriderRepository)
	})

	AfterSuite(func() {
		TruncateMenu(db)
	})

	Describe("Get By Store", func() {
		BeforeEach(func() {
			// Create Menu
			time := time.Now()
			menu, err := menuService.Create(&menu.Menu{
				Name:        "Test Menu",
				Description: "Test description menu",
				StartTime:   &time,
				EndTime:     &time,
				Enable:      false,
			})
			Expect(err).To(BeNil())
			Expect(menu).NotTo(BeNil())
		})
	})

})

func TruncateMenu(db *gorm.DB) {
	if err := db.Model(&menu.Menu{}).Delete(&menu.Menu{}).Error; err != nil {
		Expect(err).To(BeNil())
	}
}
