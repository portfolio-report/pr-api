package taxonomies

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/models"
)

type taxonomiesHandler struct {
	UserService     models.UserService
	SessionService  models.SessionService
	TaxonomyService models.TaxonomyService
	Validate        *validator.Validate
}

// NewHandler creates new taxonomies handler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	Validate *validator.Validate,
	UserService models.UserService,
	SessionService models.SessionService,
	TaxonomyService models.TaxonomyService,
) {
	h := &taxonomiesHandler{
		UserService:     UserService,
		SessionService:  SessionService,
		TaxonomyService: TaxonomyService,
		Validate:        Validate,
	}

	g := R.Group("/taxonomies")

	g.GET("/", h.GetTaxonomies)
	g.GET("/:uuid", h.GetTaxonomy)
	g.POST("/",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.PostTaxonomy)
	g.PATCH("/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.PatchTaxonomy)
	g.DELETE("/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.DeleteTaxonomy)
}
