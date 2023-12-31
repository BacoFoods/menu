package main

import (
	"fmt"
	"github.com/BacoFoods/menu/pkg/plemsi"
	"github.com/BacoFoods/menu/pkg/shared"
	"github.com/go-resty/resty/v2"
	"net/http"

	"github.com/BacoFoods/menu/pkg/connector"
	"github.com/BacoFoods/menu/pkg/scheduler"
	"github.com/BacoFoods/menu/pkg/siesa"

	"github.com/BacoFoods/menu/internal"
	"github.com/BacoFoods/menu/pkg/account"
	"github.com/BacoFoods/menu/pkg/app"
	"github.com/BacoFoods/menu/pkg/assets"
	"github.com/BacoFoods/menu/pkg/availability"
	"github.com/BacoFoods/menu/pkg/brand"
	"github.com/BacoFoods/menu/pkg/cashaudit"
	"github.com/BacoFoods/menu/pkg/category"
	"github.com/BacoFoods/menu/pkg/channel"
	"github.com/BacoFoods/menu/pkg/client"
	"github.com/BacoFoods/menu/pkg/country"
	"github.com/BacoFoods/menu/pkg/course"
	"github.com/BacoFoods/menu/pkg/currency"
	"github.com/BacoFoods/menu/pkg/database"
	"github.com/BacoFoods/menu/pkg/discount"
	"github.com/BacoFoods/menu/pkg/facturacion"
	"github.com/BacoFoods/menu/pkg/healthcheck"
	"github.com/BacoFoods/menu/pkg/invoice"
	"github.com/BacoFoods/menu/pkg/menu"
	"github.com/BacoFoods/menu/pkg/order"
	"github.com/BacoFoods/menu/pkg/payment"
	"github.com/BacoFoods/menu/pkg/payment/paymentms"
	"github.com/BacoFoods/menu/pkg/product"
	"github.com/BacoFoods/menu/pkg/router"
	"github.com/BacoFoods/menu/pkg/shift"
	"github.com/BacoFoods/menu/pkg/store"
	"github.com/BacoFoods/menu/pkg/surcharge"
	"github.com/BacoFoods/menu/pkg/swagger"
	"github.com/BacoFoods/menu/pkg/tables"
	"github.com/BacoFoods/menu/pkg/taxes"
	"github.com/BacoFoods/menu/pkg/temporal"
	"github.com/sirupsen/logrus"
)

func main() {
	// Database
	gormFramework := database.MustNewGormFramework("")
	gormDB := gormFramework.GetDBClient()

	// DB Migrations
	gormFramework.MustMakeMigrations(
		&menu.Menu{},
		&menu.MenusCategories{},
		&category.Category{},
		&discount.Discount{},
		&surcharge.Surcharge{},
		&product.Product{},
		&product.Modifier{},
		&product.Overrider{},
		&taxes.Tax{},
		&country.Country{},
		&currency.Currency{},
		&brand.Brand{},
		&store.Store{},
		&tables.Zone{},
		&tables.Table{},
		&channel.Channel{},
		&availability.Availability{},
		&order.Order{},
		&order.OrderItem{},
		&order.OrderModifier{},
		&order.OrderType{},
		&order.OrderStatus{},
		&invoice.Invoice{},
		&invoice.Item{},
		&invoice.DiscountApplied{},
		&invoice.Surcharge{},
		&account.Account{},
		&course.Course{},
		&client.Client{},
		&payment.PaymentMethod{},
		&payment.Payment{},
		&order.Attendee{},
		&shift.Shift{},
		&tables.QR{},
		&assets.Asset{},
		&cashaudit.CashAudit{},
		&cashaudit.Income{},
		&facturacion.FacturacionConfig{},
		&invoice.Document{},
		&scheduler.Schedule{},
		&connector.Equivalence{},
		&siesa.Reference{},
		&siesa.SiesaDocument{},
		&scheduler.Holiday{},
		&invoice.Resolution{},
	)

	rabbitCh := internal.MustNewRabbitMQ(internal.Config.RabbitConfig.ComandasQueue, internal.Config.RabbitConfig.Host, internal.Config.RabbitConfig.Port)
	redisConn := internal.MustNewRedis(internal.Config.RedisConfig.Host, internal.Config.RedisConfig.Port)

	httpClient := shared.NewRestClient(resty.New())

	// Healthcheck
	healthcheckHandler := healthcheck.NewHandler()
	healthcheckRoutes := healthcheck.NewRoutes(healthcheckHandler)

	// Swagger
	swaggerRoutes := swagger.NewRoutes()

	// Channel
	channelRepository := channel.NewDBRepository(gormDB)
	channelService := channel.NewService(channelRepository)
	channelHandler := channel.NewHandler(channelService)
	channelRoutes := channel.NewRoutes(channelHandler)

	// Store
	storeRepository := store.NewDBRepository(gormDB)
	storeService := store.NewService(storeRepository, channelRepository)
	storeHandler := store.NewHandler(storeService)
	storeRoutes := store.NewRoutes(storeHandler)

	// Tables
	zoneRepository := tables.NewZoneRepository(gormDB)
	tableRepository := tables.NewTableRepository(gormDB)
	tablesService := tables.NewService(tableRepository, zoneRepository, internal.Config.OITHost)
	tablesHandler := tables.NewHandler(tablesService)
	tablesRoutes := tables.NewRoutes(tablesHandler)

	// Availability
	availabilityRepository := availability.NewDBRepository(gormDB)
	availabilityService := availability.NewService(availabilityRepository, storeRepository, channelRepository)
	availabilityHandler := availability.NewHandler(availabilityService)
	availabilityRoutes := availability.NewRoutes(availabilityHandler)

	// Product
	productRepository := product.NewDBRepository(gormDB)
	productService := product.NewService(productRepository, channelRepository)
	productHandler := product.NewHandler(productService)
	productRoutes := product.NewRoutes(productHandler)

	// Surcharge
	surchargeRepository := surcharge.NewDBRepository(gormDB)
	surchargeService := surcharge.NewService(surchargeRepository)
	surchargeHandler := surcharge.NewHandler(surchargeService)
	surchargeRoutes := surcharge.NewRoutes(surchargeHandler)

	// Discount
	discountRepository := discount.NewDBRepository(gormDB)
	discountService := discount.NewService(discountRepository)
	discountHandler := discount.NewHandler(discountService)
	discountRoutes := discount.NewRoutes(discountHandler)

	// Category
	categoryRepository := category.NewDBRepository(gormDB)
	categoryService := category.NewService(categoryRepository, productRepository)
	categoryHandler := category.NewHandler(categoryService)
	categoryRoutes := category.NewRoutes(categoryHandler)

	// Menu
	menuRepository := menu.NewDBRepository(gormDB)
	menuService := menu.NewService(menuRepository, productRepository, availabilityRepository, storeRepository, categoryRepository)
	menuHandler := menu.NewHandler(menuService)
	menuRoutes := menu.NewRoutes(menuHandler)

	// Taxes
	taxesRepository := taxes.NewDBRepository(gormDB)
	taxesService := taxes.NewService(taxesRepository)
	taxesHandler := taxes.NewHandler(taxesService)
	taxesRoutes := taxes.NewRoutes(taxesHandler)

	// Country
	countryRepository := country.NewDBRepository(gormDB)
	countryService := country.NewService(countryRepository)
	countryHandler := country.NewHandler(countryService)
	countryRoutes := country.NewRoutes(countryHandler)

	// Currency
	currencyRepository := currency.NewDBRepository(gormDB)
	currencyService := currency.NewService(currencyRepository)
	currencyHandler := currency.NewHandler(currencyService)
	currencyRoutes := currency.NewRoutes(currencyHandler)

	// Paylots API
	paylotsApi := paymentms.NewPaymentsAPI(http.DefaultClient, internal.Config.PaylotsConfig.Host)

	// Payment
	paymentRepository := payment.NewDBRepository(gormDB)
	paymentService := payment.NewService(paymentRepository, paylotsApi)
	paymentHandler := payment.NewHandler(paymentService)
	paymentRoutes := payment.NewRoutes(paymentHandler)

	// Brand
	brandRepository := brand.NewDBRepository(gormDB)
	brandService := brand.NewService(brandRepository, channelRepository, paymentRepository)
	brandHandler := brand.NewHandler(brandService)
	brandRoutes := brand.NewRoutes(brandHandler)

	// Client
	clientRepository := client.NewDBRepository(gormDB)
	clientService := client.NewService(clientRepository)
	clientHandler := client.NewHandler(clientService)
	clientRoutes := client.NewRoutes(clientHandler)

	// Facturacion
	facturacionRepository := facturacion.NewRepository(gormDB)
	facturacionService := facturacion.NewService(*facturacionRepository)
	facturacionHandler := facturacion.NewHandler(facturacionService)
	facturacionRoutes := facturacion.NewRoutes(facturacionHandler)

	// Invoice
	plemsiAdapter := plemsi.NewPlemsi(httpClient)
	invoiceRepository := invoice.NewDBRepository(gormDB)
	invoiceService := invoice.NewService(invoiceRepository, clientRepository, plemsiAdapter)
	invoiceHandler := invoice.NewHandler(invoiceService)
	invoiceRoutes := invoice.NewRoutes(invoiceHandler)

	// Account
	accountRepository := account.NewDBRepository(gormDB)
	accountService := account.NewService(accountRepository)
	accountHandler := account.NewHandler(accountService)
	accountRoutes := account.NewRoutes(accountHandler)

	// Shifts
	shiftRepository := shift.NewDBRepository(gormDB)
	shiftService := shift.NewService(shiftRepository, accountRepository)
	shiftHandler := shift.NewHandler(shiftService)
	shiftRoutes := shift.NewRoutes(shiftHandler)

	// Order
	orderRepository := order.NewDBRepository(gormDB)
	orderService := order.NewService(orderRepository,
		tableRepository,
		productRepository,
		invoiceRepository,
		accountRepository,
		shiftRepository,
		rabbitCh,
		paymentService,
		discountRepository,
		channelRepository,
		facturacionService,
		facturacionRepository,
		redisConn,
		plemsiAdapter,
		clientRepository,
	)
	orderHandler := order.NewHandler(&orderService)
	orderRoutes := order.NewRoutes(orderHandler)

	// Course
	courseRepository := course.NewDBRepository(gormDB)
	courseService := course.NewService(courseRepository)
	courseHandler := course.NewHandler(courseService)
	courseRoutes := course.NewRoutes(courseHandler)

	// Siesa

	// siesaHttpClient is a special client configure to not have a timeout
	siesaHttpClient := &http.Client{
		Timeout: 0,
	}

	siesaRepository := siesa.NewDBRepository(gormDB)
	siesaService := siesa.NewService(siesaRepository, siesaHttpClient, siesa.SiesaConfig(internal.Config.SiesaConfig))
	siesaHandler := siesa.NewHandler(siesaService)
	siesaRoutes := siesa.NewRoutes(siesaHandler)

	// Temporal
	temporalHandler := temporal.NewHandler()
	temporalRoutes := temporal.NewRoutes(temporalHandler)

	// CashAudit
	cashAuditRepository := cashaudit.NewDBRepository(gormDB)
	cashAuditService := cashaudit.NewService(cashAuditRepository, storeRepository, orderRepository, invoiceRepository, shiftRepository)
	cashAuditHandler := cashaudit.NewHandler(cashAuditService)
	cashAuditRoutes := cashaudit.NewRoutes(cashAuditHandler)

	// Assets
	assetsRepository := assets.NewAssetRepository(gormDB)
	assetsService := assets.NewAssetService(assetsRepository)
	assetsHandler := assets.NewHandler(assetsService)
	assetsRoutes := assets.NewRoutes(assetsHandler)

	// Schedules
	scheduleRepository := scheduler.NewDBRepository(gormDB)
	scheduleService := scheduler.NewService(scheduleRepository)
	scheduleHandler := scheduler.NewHandler(scheduleService)
	scheduleRoutes := scheduler.NewRoutes(scheduleHandler)

	// Equivalence
	equivalenceRepository := connector.NewDBRepository(gormDB)
	equivalenceService := connector.NewService(equivalenceRepository, invoiceRepository, storeRepository)
	equivalenceHandler := connector.NewHandler(equivalenceService)
	equivalenceRoutes := connector.NewRoutes(equivalenceHandler)

	// Application TDP compilation apk and exe
	appService := app.NewService(internal.Config.GitToken, internal.Config.GitRepository)
	appHandler := app.NewHandler(appService)
	appRoutes := app.NewRoutes(appHandler)

	// Routes
	routes := &router.RoutesGroup{
		HealthCheck:  healthcheckRoutes,
		Swagger:      swaggerRoutes,
		Menu:         menuRoutes,
		Category:     categoryRoutes,
		Product:      productRoutes,
		Surcharge:    surchargeRoutes,
		Discount:     discountRoutes,
		Taxes:        taxesRoutes,
		Country:      countryRoutes,
		Currency:     currencyRoutes,
		Brand:        brandRoutes,
		Store:        storeRoutes,
		Table:        tablesRoutes,
		Channel:      channelRoutes,
		Availability: availabilityRoutes,
		Order:        orderRoutes,
		Invoice:      invoiceRoutes,
		Account:      accountRoutes,
		Course:       courseRoutes,
		Client:       clientRoutes,
		Payment:      paymentRoutes,
		Cashier:      shiftRoutes,
		Temporal:     temporalRoutes,
		CashAudit:    cashAuditRoutes,
		Assets:       assetsRoutes,
		Facturacion:  facturacionRoutes,
		Schedule:     scheduleRoutes,
		Equivalence:  equivalenceRoutes,
		Siesa:        siesaRoutes,
		App:          appRoutes,
	}

	// Run server
	r := router.NewRouter(routes)
	logrus.Fatal(r.Run(fmt.Sprintf(":%v", internal.Config.AppPort)))
}
