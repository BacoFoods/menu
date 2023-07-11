package menu_test

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
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/overriders"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/taxes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("Service Menu", func() {
	var (
		db                     *gorm.DB
		menuService            menu.Service
		availabilityRepository availability.Repository
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
		availabilityRepository = availability.NewDBRepository(db)
		storeRepository := store.NewDBRepository(db)
		categoryRepository := category.NewDBRepository(db)
		menuService = menu.NewService(menuRepository, overriderRepository, availabilityRepository, storeRepository, categoryRepository)
	})

	AfterSuite(func() {
		TruncateMenu(db)
		TruncateAvailability(db)
	})

	Describe("Create Menu", func() {
		var brandID uint
		var stores []uint
		BeforeEach(func() {
			brandID = uint(1)
			stores = []uint{101, 102, 103}
		})

		Context("When the menu is created", func() {
			It("Should return a menu object", func() {
				menu, err := menuService.Create("Menu 1", brandID, string(availability.PlaceStore), stores)
				Expect(err).To(BeNil())
				Expect(menu).NotTo(BeNil())
				Expect(menu.Name).To(Equal("Menu 1"))
				Expect(menu.BrandID).To(Equal(&brandID))
				Expect(menu.Enable).To(Equal(true))
			})

			It("Creates a store availability for each store", func() {
				menu, err := menuService.Create("Menu 2", brandID, string(availability.PlaceStore), stores)
				Expect(err).To(BeNil())
				Expect(menu).NotTo(BeNil())

				availabilities, err := availabilityRepository.FindPlacesByEntity(availability.EntityMenu, menu.ID, string(availability.PlaceStore))
				Expect(err).To(BeNil())
				Expect(availabilities).NotTo(BeNil())
				Expect(availabilities).To(HaveLen(len(stores)))

				storesAvailable := make([]uint, 0)
				for _, availability := range availabilities {
					storesAvailable = append(storesAvailable, *availability.PlaceID)
				}
				Expect(storesAvailable).To(ContainElements(stores))
			})
		})
	})
})

func TruncateMenu(db *gorm.DB) {
	if err := db.Model(&menu.Menu{}).Exec("TRUNCATE TABLE menus CASCADE").Error; err != nil {
		Expect(err).To(BeNil())
	}
}

func TruncateAvailability(db *gorm.DB) {
	if err := db.Model(&availability.Availability{}).Exec("TRUNCATE TABLE availabilities CASCADE").Error; err != nil {
		Expect(err).To(BeNil())
	}
}
