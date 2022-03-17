package portfolios

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

type PortfoliosHandler struct {
	DB               *gorm.DB
	SessionService   models.SessionService
	UserService      models.UserService
	PortfolioService models.PortfolioService
	Validate         *validator.Validate
}

func NewHandler(
	R *gin.RouterGroup,
	DB *gorm.DB,
	Validate *validator.Validate,
	SessionService models.SessionService,
	UserService models.UserService,
	PortfolioService models.PortfolioService,
) {
	h := &PortfoliosHandler{
		DB:               DB,
		SessionService:   SessionService,
		UserService:      UserService,
		PortfolioService: PortfolioService,
		Validate:         Validate,
	}

	g := R.Group("/portfolios")

	// portfolios
	g.GET("/",
		middleware.RequireUser(SessionService, UserService),
		h.GetPortfolios)
	g.POST("/",
		middleware.RequireUser(SessionService, UserService),
		h.PostPortfolios)
	g.GET("/:portfolioId",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.GetPortfolio)
	g.PUT("/:portfolioId",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.PutPortfolio)
	g.DELETE("/:portfolioId",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.DeletePortfolio)

	// securities
	g.GET("/:portfolioId/securities/",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.GetSecurities)
	g.PUT("/:portfolioId/securities/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.PutSecurity)
	g.DELETE("/:portfolioId/securities/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.DeleteSecurity)

	// accounts
	g.GET("/:portfolioId/accounts/",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.GetAccounts)
	g.PUT("/:portfolioId/accounts/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.PutAccount)
	g.DELETE("/:portfolioId/accounts/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.DeleteAccount)

	// transactions
	g.GET("/:portfolioId/transactions/",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.GetTransactions)
	g.PUT("/:portfolioId/transactions/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.PutTransaction)
	g.DELETE("/:portfolioId/transactions/:uuid",
		middleware.RequireUser(SessionService, UserService),
		middleware.RequirePortfolioPerm(PortfolioService),
		h.DeleteTransaction)
}
