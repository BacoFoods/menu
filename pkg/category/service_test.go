package category_test

import (
	"fmt"

	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/availability"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/database"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/store"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("Category", func() {
	var (
		srv     category.Service
		menuSrv menu.Service
		db      *gorm.DB
	)

	BeforeSuite(func() {
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
		)

		// Service Implementation

		categoryRepository := category.NewDBRepository(db)
		productRepository := product.NewDBRepository(db)
		menuRepository := menu.NewDBRepository(db)
		overriderRepository := product.NewDBRepository(db)
		availabilityRepository := availability.NewDBRepository(db)
		srv = category.NewService(categoryRepository, productRepository)
		menuSrv = menu.NewService(menuRepository, overriderRepository, availabilityRepository, store.NewDBRepository(db), categoryRepository)
	})

	BeforeEach(func() {
		TruncateMenu(db)
	})

	Describe("Color", func() {
		var (
			men *menu.Menu
			cat *category.Category
		)

		BeforeEach(func() {
			one := uint(1)
			var err error
			cat, err = srv.Create(&category.Category{
				Name:    "Category 1",
				Color:   "#000001",
				BrandID: &one,
			})
			Expect(err).To(BeNil())

			men, err = menuSrv.Create("Menu 1", 1, "store", []uint{1})
			Expect(err).To(BeNil())

			_, err = menuSrv.AddCategory(fmt.Sprint(men.ID), fmt.Sprint(cat.ID))
			Expect(err).To(BeNil())
		})

		It("should return the color by id", func() {
			ncat, err := srv.Get(fmt.Sprint(cat.ID))
			Expect(err).To(BeNil())
			Expect(ncat.Color).To(Equal(cat.Color))
		})

		It("should list the color from menus", func() {
			mens, err := menuSrv.FindByPlace("store", "1")
			Expect(err).To(BeNil())
			Expect(mens).ToNot(BeEmpty())
			Expect(mens[0].Categories).ToNot(BeEmpty())
			Expect(mens[0].Categories[0].Color).To(Equal("#000001"))
		})
	})
})
