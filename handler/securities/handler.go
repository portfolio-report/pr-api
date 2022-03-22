package securities

import (
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

	// public:
	g.GET("/search/:searchTerm", h.SearchSecurities)
	g.GET("/uuid/:uuid", h.GetSecurityPublic)
	g.GET("/uuid/:uuid/markets/XETR", h.GetSecurityPrices)

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

}
