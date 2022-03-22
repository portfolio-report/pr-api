package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
)

// PutAccount creates or updates account in portfolio
func (h *portfoliosHandler) PutAccount(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var req model.PortfolioAccountInput
	if err := c.BindJSON(&req); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	account, err := h.PortfolioService.UpsertPortfolioAccount(portfolioId, uuid, req)
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
	}

	c.JSON(http.StatusOK, account)
}
