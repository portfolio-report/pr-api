package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
)

// PutPortfolio updates portfolio
func (h *PortfoliosHandler) PutPortfolio(c *gin.Context) {
	portfolio := middleware.PortfolioFromContext(c)

	var input model.PortfolioInput
	if err := c.BindJSON(&input); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	portfolio, err := h.PortfolioService.UpdatePortfolio(uint(portfolio.ID), &input)
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, portfolio)
}
