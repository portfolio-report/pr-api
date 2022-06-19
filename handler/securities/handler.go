package securities

import (
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"gorm.io/gorm"
)

type securitiesHandler struct {
	*gorm.DB
	model.SessionService
	model.UserService
	*validator.Validate
	model.SecurityService
}

// NewHandler creates new securities handler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	DB *gorm.DB,
	Validate *validator.Validate,
	cacheMaxAge time.Duration,
	UserService model.UserService,
	SecurityService model.SecurityService,
	SessionService model.SessionService,
) {
	h := &securitiesHandler{
		DB:              DB,
		SessionService:  SessionService,
		SecurityService: SecurityService,
		UserService:     UserService,
		Validate:        Validate,
	}

	g := R.Group("/securities")

	memoryStore := persist.NewMemoryStore(cacheMaxAge)

	var cacheMiddleware gin.HandlerFunc
	if cacheMaxAge != 0 {
		cacheMiddleware = cache.CacheByRequestURI(memoryStore, cacheMaxAge)
	} else {
		cacheMiddleware = func(c *gin.Context) {}
	}

	// public:
	g.GET("/search/:searchTerm", h.SearchSecurities)
	g.GET("/uuid/:uuid", cacheMiddleware, h.GetSecurityPublic)
	g.GET("/uuid/:uuid/markets/:marketCode", cacheMiddleware, h.GetSecurityPrices)

	// admin:
	g.GET("/",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.GetSecurities)
	g.GET("/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.GetSecurityAdmin)
	g.POST("/",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.PostSecurity)
	g.PATCH("/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.PatchSecurity)
	g.DELETE("/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.DeleteSecurity)
	g.POST("/uuid/:uuid/logo",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.UpdateLogo)
	g.PATCH("/uuid/:uuid/markets/:marketCode",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.PatchSecurityMarket)
	g.DELETE("/uuid/:uuid/markets/:marketCode",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.DeleteSecurityMarket)
	g.PUT("/uuid/:uuid/taxonomies/:rootUuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.PutSecurityTaxonomies)
	g.GET("/maintenance/gaps",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.GetGaps)

}
