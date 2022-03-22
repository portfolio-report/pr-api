package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/graph/model"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.

// Resolver contains dependencies to be injected
type Resolver struct {
	*gorm.DB
	*validator.Validate
	model.UserService
	model.SessionService
	model.PortfolioService
	model.CurrenciesService
	model.SecurityService
}
