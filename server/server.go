package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/service"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func setupCron(cs model.CurrenciesService) {
	logger := log.New(os.Stderr, "[cron] ", log.LstdFlags|log.Lmsgprefix)

	updateExchangeRates := func() {
		// Recover from panic
		defer func() {
			if r := recover(); r != nil {
				logger.Println("Panic while updating exchange rates:", r,
					"\nstacktrace:\n"+string(debug.Stack()))
			}
		}()

		if err := cs.UpdateExchangeRates(); err != nil {
			logger.Println("Error while updating exchange rates:", err)
		}
	}

	go func() {
		// Run once after 5min, then every 2hours
		time.Sleep(5 * time.Minute)
		for {
			updateExchangeRates()
			time.Sleep(2 * time.Hour)
		}
	}()
}

func PrepareApp() (*service.Config, *gorm.DB) {
	// Read config
	cfg := service.ReadConfig()

	// Default to Gin Release Mode
	if os.Getenv(gin.EnvGinMode) == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Fatalln("Could not connect to database:", err)
	}

	return cfg, db
}

func InitializeService(cfg *service.Config, db *gorm.DB) *handler.Config {
	// Initialize validator
	validate := libs.GetValidator()

	// Initialize services
	currenciesService := service.NewCurrenciesService(db, true)
	geoipService := service.NewGeoipService(cfg.Ip2locToken)
	userService := service.NewUserService(db)
	sessionService := service.NewSessionService(db, validate, cfg.SessionTimeout)
	portfolioService := service.NewPortfolioService(db)
	securityService := service.NewSecurityService(db)
	taxonomyService := service.NewTaxonomyService(db, validate)
	mailerService, err := service.NewMailerService(cfg.MailerTransport, cfg.ContactRecipientEmail, validate)
	if err != nil {
		fmt.Println("WARNING: Cannot send emails, could not create MailerService: " + err.Error())
	}

	// Setup cronjobs
	setupCron(currenciesService)

	// Register custom validations on GIN validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		libs.RegisterCustomValidations(v)
	}

	return &handler.Config{
		UserService:       userService,
		SessionService:    sessionService,
		CurrenciesService: currenciesService,
		MailerService:     mailerService,
		GeoipService:      geoipService,
		PortfolioService:  portfolioService,
		SecurityService:   securityService,
		TaxonomyService:   taxonomyService,
		BaseURL:           "",
		CacheMaxAge:       cfg.CacheMaxAge,
		DB:                db,
		Validate:          validate,
	}
}

func CreateApp(cfg *handler.Config) http.Handler {
	router := gin.New()

	router.Use(gin.Logger())

	router.NoRoute(func(c *gin.Context) { libs.HandleNotFoundError(c) })

	// Recover from any panics, log stack trace and reply with internal server error
	router.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		code := http.StatusInternalServerError
		c.JSON(code, gin.H{"statusCode": code, "error": http.StatusText(code)})
	}))

	// Use CORS middleware
	router.Use(libs.Cors)

	handler.NewHandler(router, cfg)

	return router
}
