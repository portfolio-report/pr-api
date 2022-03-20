package models

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/shopspring/decimal"
)

// CurrenciesService describes the interface of currencies service
type CurrenciesService interface {
	GetCurrencies() ([]*model.Currency, error)
	GetExchangerate(baseCC, quoteCC string) (*model.Exchangerate, error)
	GetExchangeratePrices(er *model.Exchangerate, from *string) ([]*model.ExchangeratePrice, error)
	ConvertCurrencyAmount(decimal.Decimal, string, string, time.Time) (decimal.Decimal, error)
	UpdateExchangeRates() error
}

// GeoipService describes the interface of GeoIP service
type GeoipService interface {
	GetCountryFromIp(string) string
}

// MailerService describes the interface of mailer service
type MailerService interface {
	SendContactMail(senderEmail string, senderName string, subject string, message string, ip string) error
}

// PortfolioService describes the interface of portfolio service
type PortfolioService interface {
	GetPortfolioByID(ID uint) (*model.Portfolio, error)
	GetPortfolioOfUserByID(user *model.User, ID uint) (*model.Portfolio, error)
	GetAllOfUser(user *model.User) ([]*model.Portfolio, error)
	CreatePortfolio(user *model.User, req *model.PortfolioInput) (*model.Portfolio, error)
	UpdatePortfolio(ID uint, req *model.PortfolioInput) (*model.Portfolio, error)
	DeletePortfolio(ID uint) (*model.Portfolio, error)
}

// SecurityService describes the interface of security service
type SecurityService interface {
	GetSecurityByUUID(uuid string) (*model.Security, error)
	GetEventsOfSecurity(security *model.Security) ([]*model.Event, error)
	CreateSecurity(input *model.SecurityInput) (*model.Security, error)
	UpdateSecurity(uuid string, input *model.SecurityInput) (*model.Security, error)
	DeleteSecurity(uuid string) (*model.Security, error)
}

// SessionService describes the interface of session service
type SessionService interface {
	GetAllOfUser(user *model.User) ([]*model.Session, error)
	CreateSession(user *model.User, note string) (*model.Session, error)
	DeleteSession(token string) (*model.Session, error)
	GetSessionToken(c *gin.Context) string
	ValidateToken(token string) (*model.Session, error)
	CleanupExpiredSessions() error
}

// TaxonomyService describes the interface of taxonomy service
type TaxonomyService interface {
	GetAllTaxonomies() ([]*model.Taxonomy, error)
	GetTaxonomyByUUID(uuid string) (*model.Taxonomy, error)
	GetDescendantsOfTaxonomy(taxonomy *model.Taxonomy) ([]*model.Taxonomy, error)
	CreateTaxonomy(taxonomy *model.Taxonomy) (*model.Taxonomy, error)
	UpdateTaxonomy(uuid string, taxonomy *model.Taxonomy) (*model.Taxonomy, error)
	DeleteTaxonomy(uuid string) (*model.Taxonomy, error)
}

// UserService describes the interface of user service
type UserService interface {
	Create(username string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserFromSession(session *model.Session) (*model.User, error)
	UpdatePassword(ctx context.Context, user *model.User, password string) error
	VerifyPassword(ctx context.Context, user *model.User, password string) (bool, error)
	Delete(user *model.User) error
	UpdateLastSeen(user *model.User) error
}
