package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
)

// PutTransaction creates or updates transaction in portfolio
func (h *portfoliosHandler) PutTransaction(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var req model.PortfolioTransactionInput
	if err := c.BindJSON(&req); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	transaction, err := h.PortfolioService.UpsertPortfolioTransaction(portfolioId, uuid, req)
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, transaction)
}
