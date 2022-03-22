package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

// GetTransactions lists all transactions in portfolio
func (h *portfoliosHandler) GetTransactions(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID

	transactions, _ := h.PortfolioService.GetPortfolioTransactionsOfPortfolio(portfolioId)

	c.JSON(http.StatusOK, transactions)
}
