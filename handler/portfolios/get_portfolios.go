package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

// GetPortfolios lists all portfolios of current user
func (h *portfoliosHandler) GetPortfolios(c *gin.Context) {
	user := middleware.UserFromContext(c.Request.Context())
	portfolios := h.PortfolioService.GetAllOfUser(user)
	c.JSON(http.StatusOK, portfolios)
}
