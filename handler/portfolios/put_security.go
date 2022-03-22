package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
)

// PutSecurity creates or updates security in portfolio
func (h *portfoliosHandler) PutSecurity(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var req model.PortfolioSecurityInput
	if err := c.BindJSON(&req); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	security, err := h.PortfolioService.UpsertPortfolioSecurity(portfolioId, uuid, req)
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, security)
}
