package taxonomies

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

type taxonomiesHandler struct {
	model.UserService
	model.SessionService
	model.TaxonomyService
	*validator.Validate
}

// NewHandler creates new taxonomies handler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	Validate *validator.Validate,
	UserService model.UserService,
	SessionService model.SessionService,
	TaxonomyService model.TaxonomyService,
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
	g.PUT("/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.PutTaxonomy)
	g.DELETE("/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.DeleteTaxonomy)
}
