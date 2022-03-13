package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/models"
)

func (h *PortfoliosHandler) GetTransactions(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID

	var transactions []db.PortfolioTransaction
	err := h.DB.
		Preload("Units").
		Where("portfolio_id = ?", portfolioId).
		Find(&transactions).Error
	if err != nil {
		panic(err)
	}

	response := []models.PortfolioTransactionResponse{}
	for _, db := range transactions {
		response = append(response, models.PortfolioTransactionResponseFromDB(&db))
	}

	c.JSON(http.StatusOK, response)
}
