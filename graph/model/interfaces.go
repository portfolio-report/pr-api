package model

import (
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

// CurrenciesService describes the interface of currencies service
type CurrenciesService interface {
	GetCurrencies() []*Currency
	GetExchangerate(baseCC, quoteCC string) (*Exchangerate, error)
	GetExchangeratePrices(exchangerateID uint, from *string) ([]*ExchangeratePrice, error)
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
	GetPortfolioByID(ID uint) (*Portfolio, error)
	GetPortfolioOfUserByID(user *User, ID uint) (*Portfolio, error)
	GetAllOfUser(user *User) []*Portfolio
	CreatePortfolio(user *User, req *PortfolioInput) (*Portfolio, error)
	UpdatePortfolio(ID uint, req *PortfolioInput) (*Portfolio, error)
	DeletePortfolio(ID uint) *Portfolio

	GetPortfolioAccountsOfPortfolio(portfolioId int) []*PortfolioAccount
	UpsertPortfolioAccount(portfolioId int, uuid uuid.UUID, input PortfolioAccountInput) (*PortfolioAccount, error)
	DeletePortfolioAccount(portfolioId int, uuid uuid.UUID) (*PortfolioAccount, error)

	GetPortfolioSecuritiesOfPortfolio(portfolioId int) []*PortfolioSecurity
	UpsertPortfolioSecurity(portfolioId int, uuid uuid.UUID, input PortfolioSecurityInput) (*PortfolioSecurity, error)
	DeletePortfolioSecurity(portfolioId int, uuid uuid.UUID) (*PortfolioSecurity, error)
	CalcSecurityShares(securities []PortfolioSecurityKey) []*decimal.Decimal

	GetPortfolioTransactionsOfPortfolio(portfolioId int) []*PortfolioTransaction
	UpsertPortfolioTransaction(portfolioId int, uuid uuid.UUID, input PortfolioTransactionInput) (*PortfolioTransaction, error)
	DeletePortfolioTransaction(portfolioId int, uuid uuid.UUID) (*PortfolioTransaction, error)
}

// SecurityService describes the interface of security service
type SecurityService interface {
	GetSecurityByUUID(uuid uuid.UUID) (*Security, error)
	GetSecuritiesByTag(tag string) []*Security
	GetEventsOfSecurity(security *Security) []*Event
	CreateSecurity(input *SecurityInput) (*Security, error)
	UpdateSecurity(uuid uuid.UUID, input *SecurityInput) (*Security, error)
	DeleteSecurity(uuid uuid.UUID) (*Security, error)
	UpdateLogo(uuid uuid.UUID, logo io.Reader, extension string) (string, error)
	DeleteSecurityMarket(securityUuid uuid.UUID, marketCode string) (*SecurityMarket, error)
	UpdateSecurityTaxonomies(securityUuid, rootTaxonomyUuid uuid.UUID, inputs []*SecurityTaxonomyInput) ([]*SecurityTaxonomy, error)
	UpsertTag(name string, securityUuids []uuid.UUID) ([]*Security, error)
	DeleteTag(name string)
	FindGapsInPrices(minDuration, maxResults int) []map[string]interface{}
	LogoUrlFromExtras(extrasJson datatypes.JSON) *string
}

// SessionService describes the interface of session service
type SessionService interface {
	GetAllOfUser(user *User) []*Session
	CreateSession(user *User, note string) (*Session, error)
	DeleteSession(token string) (*Session, error)
	GetSessionToken(c *gin.Context) string
	ValidateToken(token string) (*Session, error)
	CleanupExpiredSessions() error
}

// TaxonomyService describes the interface of taxonomy service
type TaxonomyService interface {
	GetAllTaxonomies() []*Taxonomy
	GetTaxonomyByUUID(uuid uuid.UUID) (*Taxonomy, error)
	GetDescendantsOfTaxonomy(taxonomy *Taxonomy) []*Taxonomy
	CreateTaxonomy(taxonomy *TaxonomyInput) (*Taxonomy, error)
	UpdateTaxonomy(uuid uuid.UUID, taxonomy *TaxonomyInput) (*Taxonomy, error)
	DeleteTaxonomy(uuid uuid.UUID) (*Taxonomy, error)
}

// UserService describes the interface of user service
type UserService interface {
	Create(username string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetByIDs(ids []int) ([]*User, error)
	GetUserFromSession(session *Session) (*User, error)
	UpdatePassword(ctx context.Context, user *User, password string) error
	VerifyPassword(ctx context.Context, user *User, password string) (bool, error)
	Delete(id int) error
	UpdateLastSeen(user *User) error
}
