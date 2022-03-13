package securities

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

type SecuritiesHandler struct {
	DB             *gorm.DB
	SessionService models.SessionService
	UserService    models.UserService
	validate       *validator.Validate
}

func NewHandler(
	R *gin.RouterGroup,
	DB *gorm.DB,
	validate *validator.Validate,
	UserService models.UserService,
	SessionService models.SessionService,
) {
	h := &SecuritiesHandler{
		DB:             DB,
		SessionService: SessionService,
		UserService:    UserService,
		validate:       validate,
	}

	g := R.Group("/securities")

	// public:
	g.GET("/search/:searchTerm", h.SearchSecurities)
	g.GET("/uuid/:uuid", h.GetSecurityPublic)
	g.GET("/uuid/:uuid/markets/XETR", h.GetSecurityPrices)

	// admin:
	g.GET("",
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
		h.PostSecurityAdmin)
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
