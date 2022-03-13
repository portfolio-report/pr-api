package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

func (h *PortfoliosHandler) DeletePortfolio(c *gin.Context) {
	portfolio := middleware.PortfolioFromContext(c)
	h.PortfolioService.DeletePortfolio(uint(portfolio.ID))
	c.JSON(http.StatusOK, portfolio)
}
