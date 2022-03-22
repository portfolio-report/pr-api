package handler

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/auth"
	"github.com/portfolio-report/pr-api/handler/currencies"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/handler/portfolios"
	"github.com/portfolio-report/pr-api/handler/securities"
	"github.com/portfolio-report/pr-api/handler/stats"
	"github.com/portfolio-report/pr-api/handler/taxonomies"
	"gorm.io/gorm"
)

// Config holds configuration for all handlers
type Config struct {
	model.UserService
	model.SessionService
	model.CurrenciesService
	model.PortfolioService
	model.SecurityService
	model.TaxonomyService
	model.MailerService
	model.GeoipService
	BaseURL string
	*gorm.DB
	*validator.Validate
}

type rootHandler struct {
	model.UserService
	model.SessionService
	model.CurrenciesService
	model.PortfolioService
	model.SecurityService
	model.TaxonomyService
	model.MailerService
	model.GeoipService
	*gorm.DB
	*validator.Validate
}

// NewHandler creates new root handler and registers routes
func NewHandler(R *gin.Engine, c *Config) {
	h := &rootHandler{
		UserService:       c.UserService,
		SessionService:    c.SessionService,
		CurrenciesService: c.CurrenciesService,
		PortfolioService:  c.PortfolioService,
		SecurityService:   c.SecurityService,
		TaxonomyService:   c.TaxonomyService,
		MailerService:     c.MailerService,
		GeoipService:      c.GeoipService,
		DB:                c.DB,
		Validate:          c.Validate,
	}

	R.Use(middleware.AuthUser(c.SessionService, c.UserService))

	g := R.Group(c.BaseURL)

	g.GET("", h.GetRoot)

	// /graphql
	g.POST("/graphql", middleware.Useragent, h.GraphqlHandler())
	g.GET("/graphql", h.PlaygroundHandler(path.Join(g.BasePath(), "graphql")))

	// /doc
	h.RegisterSwaggerUi(g, "/doc")

	// /auth
	auth.NewHandler(g, c.DB, c.Validate, c.SessionService, c.UserService)

	// /currencies
	currencies.NewHandler(g, c.UserService, c.SessionService, c.CurrenciesService)

	// /stats
	stats.NewHandler(g, c.DB, c.UserService, c.SessionService, c.GeoipService)

	// /contact
	g.POST("/contact", h.Contact)

	// /securities
	securities.NewHandler(g, c.DB, c.Validate, c.UserService, c.SecurityService, c.SessionService)

	// /portfolios
	portfolios.NewHandler(g, c.DB, c.SessionService, c.UserService, c.PortfolioService)

	// /taxonomies
	taxonomies.NewHandler(g, c.Validate, c.UserService, c.SessionService, c.TaxonomyService)

}
