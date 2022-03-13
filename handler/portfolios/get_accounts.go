package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/models"
)

func (h *PortfoliosHandler) GetAccounts(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID

	var accounts []db.PortfolioAccount
	err := h.DB.Where("portfolio_id = ?", portfolioId).Find(&accounts).Error
	if err != nil {
		panic(err)
	}

	response := []models.PortfolioAccountResponse{}
	for _, db := range accounts {
		response = append(response, models.PortfolioAccountResponseFromDB(&db))
	}

	c.JSON(http.StatusOK, response)
}
