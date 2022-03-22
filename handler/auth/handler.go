package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"gorm.io/gorm"
)

type authHandler struct {
	*gorm.DB
	model.SessionService
	model.UserService
	*validator.Validate
}

// NewHandler creates new AuthHandler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	DB *gorm.DB,
	Validate *validator.Validate,
	SessionService model.SessionService,
	UserService model.UserService,
) {
	h := &authHandler{
		DB:             DB,
		SessionService: SessionService,
		UserService:    UserService,
		Validate:       Validate,
	}

	g := R.Group("/auth")

	g.POST("/register", middleware.Useragent, h.RegisterUser)
	g.POST("/login", middleware.Useragent, h.LoginUser)
	g.POST("/logout", middleware.RequireUser(SessionService, UserService), h.LogoutUser)

	g.GET("/users/me", middleware.RequireUser(SessionService, UserService), h.GetMe)
	g.DELETE("/users/me", middleware.RequireUser(SessionService, UserService), h.DeleteMe)
	g.POST("/users/me/password", middleware.RequireUser(SessionService, UserService), h.UpdatePassword)

	g.GET("/sessions", middleware.RequireUser(SessionService, UserService), h.GetSessions)
}
