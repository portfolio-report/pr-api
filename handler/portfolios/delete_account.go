package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
)

// DeleteAccount removes account from portfolio and links to it
func (h *portfoliosHandler) DeleteAccount(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	account, err := h.PortfolioService.DeletePortfolioAccount(portfolioId, uuid)
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	c.JSON(http.StatusOK, account)
}
