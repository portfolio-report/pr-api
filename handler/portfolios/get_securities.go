package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/models"
)

// GetSecurities lists all securities in portfolio
func (h *PortfoliosHandler) GetSecurities(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID

	var securities []db.PortfolioSecurity
	err := h.DB.Where("portfolio_id = ?", portfolioId).Find(&securities).Error
	if err != nil {
		panic(err)
	}

	response := []models.PortfolioSecurityResponse{}
	for _, db := range securities {
		response = append(response, models.PortfolioSecurityResponseFromDB(&db))
	}

	c.JSON(http.StatusOK, response)
}
