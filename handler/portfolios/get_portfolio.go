package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

// GetPortfolio returns single portfolio
func (h *portfoliosHandler) GetPortfolio(c *gin.Context) {
	portfolio := middleware.PortfolioFromContext(c)
	c.JSON(http.StatusOK, portfolio)
}
