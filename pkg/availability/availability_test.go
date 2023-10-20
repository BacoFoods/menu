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
		storeService        store.Service
		brandService        brand.Service
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

		// Service Implementation
		brandRepository := brand.NewDBRepository(db)
		brandService = brand.NewService(brandRepository)

		storeRepository := store.NewDBRepository(db)
		channelRepository := channel.NewDBRepository(db)
		storeService = store.NewService(storeRepository, channelRepository)

		overridersRepository := overriders.NewDBRepository(db)
		availabilityRepository := availability.NewDBRepository(db)
		availabilityService = availability.NewService(availabilityRepository, storeRepository, channelRepository)

		menuRepository := menu.NewDBRepository(db)
		categoryRepository := category.NewDBRepository(db)
		menuService = menu.NewService(menuRepository, overridersRepository, availabilityRepository, storeRepository, categoryRepository)
	})

	AfterSuite(func() {
		TruncateAvailability(db)
		TruncateMenu(db)
	})

	BeforeEach(func() {
		TruncateAvailability(db)
		TruncateMenu(db)
	})

	var (
		myBrand *brand.Brand
		myStore *store.Store
	)

	BeforeEach(func() {
		var err error
		myBrand, err = brandService.Create(&brand.Brand{
			Name: "brand test",
		})
		Expect(err).To(BeNil())
		Expect(myBrand).To(Not(BeNil()))

		myStore, err = storeService.Create(&store.Store{
			Name:    "store test",
			BrandID: &myBrand.ID,
			Enabled: true,
		})
		Expect(err).To(BeNil())
		Expect(myStore).To(Not(BeNil()))
		Expect(myStore.ID).To(Not(BeNil()))
	})

	Describe("Remove Entity", func() {
		Context("Entity doesnt exists", func() {
			It("Should not return error", func() {
				// Arrange
				place := availability.PlaceStore
				placeID := myStore.ID
				entityID := uint(10)

				// Act
				err := availabilityService.RemoveEntity(availability.EntityMenu, place, entityID, placeID)

				// Assert
				Expect(err).To(BeNil())
			})
		})

		Context("Entity exists", func() {
			var (
				place    = availability.PlaceStore
				placeID  uint
				entityID = uint(11)
			)

			BeforeEach(func() {
				placeID = myStore.ID
				err := availabilityService.EnableEntity(availability.EntityMenu, place, entityID, placeID, true)
				Expect(err).To(BeNil())

				av, err := availabilityService.Get(availability.EntityMenu, place, entityID, placeID)
				Expect(err).To(BeNil())
				Expect(av).ToNot(BeNil())
			})

			It("Should not return error", func() {
				// Act
				err := availabilityService.RemoveEntity(availability.EntityMenu, place, entityID, placeID)

				// Assert
				Expect(err).To(BeNil())
			})

			It("Should remove the entity", func() {
				// Act
				err := availabilityService.RemoveEntity(availability.EntityMenu, place, entityID, placeID)

				// Assert
				Expect(err).To(BeNil())

				_, err = availabilityService.Get(availability.EntityMenu, place, entityID, placeID)
				Expect(err).To(Equal(gorm.ErrRecordNotFound))
			})
		})
	})

	Describe("Enable/Disable Entity", func() {
		Context("When the entity is a menu", func() {
			It("Should enable/disable the entity", func() {
				// Arrange
				place := availability.PlaceStore
				placeIDs := []uint{myStore.ID}
				menu, err := menuService.Create("menu test", myBrand.ID, string(place), placeIDs)
				Expect(err).To(BeNil())
				Expect(menu).To(Not(BeNil()))

				// Act
				err = availabilityService.EnableEntity(availability.EntityMenu, place, menu.ID, placeIDs[0], false)
				availability, err := availabilityService.Get(availability.EntityMenu, availability.PlaceStore, menu.ID, placeIDs[0])

				// Assert
				Expect(err).To(BeNil())
				Expect(availability).To(Not(BeNil()))
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
