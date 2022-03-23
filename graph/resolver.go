package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/portfolio-report/pr-api/graph/model"
)

// This file will not be regenerated automatically.

// Resolver contains dependencies to be injected
type Resolver struct {
	model.UserService
	model.SessionService
	model.PortfolioService
	model.CurrenciesService
	model.SecurityService
}
