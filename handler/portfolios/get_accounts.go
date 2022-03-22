package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

// GetAccounts lists all accounts in portfolio
func (h *portfoliosHandler) GetAccounts(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID

	accounts, _ := h.PortfolioService.GetPortfolioAccountsOfPortfolio(portfolioId)

	c.JSON(http.StatusOK, accounts)
}
