package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB             *gorm.DB
	SessionService models.SessionService
	UserService    models.UserService
	Validate       *validator.Validate
}

func NewHandler(
	R *gin.RouterGroup,
	DB *gorm.DB,
	Validate *validator.Validate,
	SessionService models.SessionService,
	UserService models.UserService,
) {
	h := &AuthHandler{
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
