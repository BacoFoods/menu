package store_test

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
)

var _ = Describe("Store Service", func() {

	var (
		db           *gorm.DB
		storeService store.Service
		brandService brand.Service
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
		storeRepository := store.NewDBRepository(db)
		storeService = store.NewService(storeRepository)

		// Brand Implementation
		brandRepository := brand.NewDBRepository(db)
		brandService = brand.NewService(brandRepository)
	})

	AfterSuite(func() {
		TruncateStore(db)
		TruncateBrand(db)
	})

	Describe("Get By Brand", func() {
		var brandID *uint
		var storeID uint

		BeforeEach(func() {
			// Create brand
			brand, err := brandService.Create(&brand.Brand{
				Name: "Test Brand",
			})
			Expect(err).To(BeNil())
			Expect(brand).ToNot(BeNil())
			brandID = &brand.ID

			// Create store
			store, err := storeService.Create(&store.Store{
				BrandID: &brand.ID,
				Name:    "Test Store",
			})
			Expect(err).To(BeNil())
			Expect(store).ToNot(BeNil())
			storeID = store.ID
		})

		Context("When the brand does not exist", func() {
			It("Should return an empty list", func() {
				// Given
				TruncateBrand(db)

				filter := make(map[string]string)
				filter["brand_id"] = fmt.Sprintf("%d", *brandID)

				// When
				stores, err := storeService.Find(filter)

				// Then
				Expect(err).To(BeNil())
				Expect(stores).To(HaveLen(0))
			})
		})

		Context("When the brand exists", func() {
			It("Should return the store", func() {
				// When
				filter := make(map[string]string)
				filter["brand_id"] = fmt.Sprintf("%d", *brandID)
				stores, err := storeService.Find(filter)

				// Then
				Expect(err).To(BeNil())
				Expect(stores).ToNot(BeNil())
				Expect(stores[0].BrandID).To(Equal(brandID))
				Expect(stores[0].ID).To(Equal(storeID))
			})
		})

		Context("When the store does not exist", func() {
			It("Should return an empty list", func() {
				// Given
				TruncateStore(db)

				// When
				filter := make(map[string]string)
				filter["brand_id"] = fmt.Sprintf("%d", *brandID)
				stores, err := storeService.Find(filter)

				// Then
				Expect(err).To(BeNil())
				Expect(stores).To(HaveLen(0))
			})
		})
	})
})

func TruncateBrand(db *gorm.DB) {
	if err := db.Model(&brand.Brand{}).Exec("TRUNCATE TABLE brands CASCADE").Error; err != nil {
		Expect(err).To(BeNil())
	}
}

func TruncateStore(db *gorm.DB) {
	if err := db.Model(&store.Store{}).Exec("TRUNCATE TABLE stores CASCADE").Error; err != nil {
		Expect(err).To(BeNil())
	}
}
