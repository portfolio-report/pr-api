package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/handler"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"github.com/portfolio-report/pr-api/service"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func setupCron(cs models.CurrenciesService) {
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
		for true {
			updateExchangeRates()
			time.Sleep(2 * time.Hour)
		}
	}()
}

func main() {
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

	// Initialize validator
	validate := libs.GetValidator()

	// Initialize services
	currenciesService := service.NewCurrenciesService(db)
	geoipService := service.NewGeoipService(cfg.Ip2locToken)
	userService := service.NewUserService(db)
	sessionService := service.NewSessionService(db, validate, cfg.SessionTimeout)
	portfolioService := service.NewPortfolioService(db)
	securityService := service.NewSecurityService(db)
	taxonomyService := service.NewTaxonomyService(db)
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

	handler.NewHandler(&handler.Config{
		R:                 router,
		UserService:       userService,
		SessionService:    sessionService,
		CurrenciesService: currenciesService,
		MailerService:     mailerService,
		GeoipService:      geoipService,
		PortfolioService:  portfolioService,
		SecurityService:   securityService,
		TaxonomyService:   taxonomyService,
		BaseURL:           "",
		DB:                db,
		Validate:          validate,
	})

	address := ":3000"

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Println("Listening and serving HTTP on", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Received signal to stop, shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutdown gracefully: ", err)
	}

	log.Println("Server terminated.")
}
