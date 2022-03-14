package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/handler/middleware"
)

func (h *PortfoliosHandler) GetPortfolios(c *gin.Context) {
	user := middleware.UserFromContext(c.Request.Context())
	portfolios, err := h.PortfolioService.GetAllOfUser(user)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, portfolios)
}
