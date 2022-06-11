package tags

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

type tagsHandler struct {
	model.UserService
	model.SessionService
	model.SecurityService
	*validator.Validate
}

// NewHandler creates new tags handler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	Validate *validator.Validate,
	UserService model.UserService,
	SessionService model.SessionService,
	SecurityService model.SecurityService,
) {
	h := &tagsHandler{
		UserService:     UserService,
		SecurityService: SecurityService,
		SessionService:  SessionService,
		Validate:        Validate,
	}

	g := R.Group("/tags")

	g.GET("/:name", h.GetTag)
	g.PUT("/:name",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.PutTag)
	g.DELETE("/:name",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequireAdmin(),
		h.DeleteTag)
}
