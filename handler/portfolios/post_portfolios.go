package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
)

func (h *PortfoliosHandler) PostPortfolios(c *gin.Context) {
	user := middleware.UserFromContext(c.Request.Context())

	var input model.PortfolioInput
	if err := c.BindJSON(&input); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	portfolio, err := h.PortfolioService.CreatePortfolio(user, &input)
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, portfolio)
}
