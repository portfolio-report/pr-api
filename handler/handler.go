package handler

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/handler/auth"
	"github.com/portfolio-report/pr-api/handler/currencies"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/handler/portfolios"
	"github.com/portfolio-report/pr-api/handler/securities"
	"github.com/portfolio-report/pr-api/handler/stats"
	"github.com/portfolio-report/pr-api/handler/taxonomies"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

type Config struct {
	R                 *gin.Engine
	UserService       models.UserService
	SessionService    models.SessionService
	CurrenciesService models.CurrenciesService
	PortfolioService  models.PortfolioService
	SecurityService   models.SecurityService
	TaxonomyService   models.TaxonomyService
	MailerService     models.MailerService
	GeoipService      models.GeoipService
	BaseURL           string
	DB                *gorm.DB
	Validate          *validator.Validate
}

type Handler struct {
	UserService       models.UserService
	SessionService    models.SessionService
	CurrenciesService models.CurrenciesService
	PortfolioService  models.PortfolioService
	SecurityService   models.SecurityService
	TaxonomyService   models.TaxonomyService
	MailerService     models.MailerService
	GeoipService      models.GeoipService
	DB                *gorm.DB
	validate          *validator.Validate
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:       c.UserService,
		SessionService:    c.SessionService,
		CurrenciesService: c.CurrenciesService,
		PortfolioService:  c.PortfolioService,
		SecurityService:   c.SecurityService,
		TaxonomyService:   c.TaxonomyService,
		MailerService:     c.MailerService,
		GeoipService:      c.GeoipService,
		DB:                c.DB,
		validate:          c.Validate,
	}

	c.R.Use(middleware.AuthUser(c.SessionService, c.UserService))

	g := c.R.Group(c.BaseURL)

	g.GET("", h.GetRoot)

	// /graphql
	g.POST("/graphql", middleware.Useragent, h.GraphqlHandler())
	g.GET("/graphql", h.PlaygroundHandler(path.Join(g.BasePath(), "graphql")))

	// /doc
	h.RegisterSwaggerUi(g, "/doc")

	// /auth
	auth.NewHandler(g, c.DB, c.Validate, c.SessionService, c.UserService)

	// /currencies
	currencies.NewHandler(g, c.DB, c.UserService, c.SessionService, c.CurrenciesService)

	// /stats
	stats.NewHandler(g, c.DB, c.UserService, c.SessionService, c.GeoipService)

	// /contact
	g.POST("/contact", h.Contact)

	// /securities
	securities.NewHandler(g, c.DB, c.Validate, c.UserService, c.SessionService)

	// /portfolios
	portfolios.NewHandler(g, c.DB, c.Validate, c.SessionService, c.UserService, c.PortfolioService)

	// /taxonomies
	taxonomies.NewHandler(g, c.DB, c.Validate, c.UserService, c.SessionService, c.TaxonomyService)

}
