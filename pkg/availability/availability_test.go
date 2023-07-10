package availability_test

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

	"github.com/BacoFoods/menu/pkg/availability"
)

var _ = Describe("Availability", func() {
	var (
		db                  *gorm.DB
		menuService         menu.Service
		availabilityService availability.Service
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
		overridersRepository := overriders.NewDBRepository(db)
		availabilityRepository := availability.NewDBRepository(db)
		availabilityService = availability.NewService(availabilityRepository)

		menuRepository := menu.NewDBRepository(db)
		menuService = menu.NewService(menuRepository, overridersRepository, availabilityRepository)
	})

	AfterSuite(func() {
		TruncateAvailability(db)
		TruncateMenu(db)
	})

	Describe("Enable/Disable Entity", func() {
		Context("When the entity is a menu", func() {
			It("Should enable/disable the entity", func() {
				// Arrange
				place := availability.PlaceStore
				placeIDs := []uint{4, 2, 3}
				brandID := uint(1)
				menu, err := menuService.Create("menu test", brandID, string(place), placeIDs)
				Expect(err).To(BeNil())
				Expect(menu).To(Not(BeNil()))
				Expect(menu).To(Not(BeNil()))

				// Act
				err = availabilityService.EnableEntity(availability.EntityMenu, place, placeIDs[0], menu.ID, false)
				availability, err := availabilityService.Get(availability.EntityMenu, availability.PlaceStore, placeIDs[0], menu.ID)

				// Assert
				Expect(err).To(BeNil())
				Expect(availability).To(Not(BeNil()))
				Expect(availability.Enable).To(BeFalse())
			})
		})
	})
})

func TruncateAvailability(db *gorm.DB) {
	if err := db.Model(&availability.Availability{}).Exec("TRUNCATE TABLE availabilities CASCADE").Error; err != nil {
		Expect(err).To(BeNil())
	}
}

func TruncateMenu(db *gorm.DB) {
	if err := db.Model(&menu.Menu{}).Exec("TRUNCATE TABLE menus CASCADE").Error; err != nil {
		Expect(err).To(BeNil())
	}
}
