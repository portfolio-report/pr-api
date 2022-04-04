package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

// GetSecurities lists all securities in portfolio
func (h *portfoliosHandler) GetSecurities(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID
	securities := h.PortfolioService.GetPortfolioSecuritiesOfPortfolio(portfolioId)
	c.JSON(http.StatusOK, securities)
}
