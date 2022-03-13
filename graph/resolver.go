package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB               *gorm.DB
	Validate         *validator.Validate
	UserService      models.UserService
	SessionService   models.SessionService
	PortfolioService models.PortfolioService
	SecurityService  models.SecurityService
}
