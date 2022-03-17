package portfolios

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm/clause"
)

// DeleteTransaction removes transaction from portfolio and links to it
func (h *portfoliosHandler) DeleteTransaction(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	// Remove link from partner transaction (if exists)
	h.DB.Model(&db.PortfolioTransaction{}).
		Where("portfolio_id = ? AND partner_transaction_uuid = ?", portfolioId, uuid).
		Update("partner_transaction_uuid", nil)

	var transaction db.PortfolioTransaction
	result := h.DB.Clauses(clause.Returning{}).Where("portfolio_id = ? AND uuid = ?", portfolioId, uuid).Delete(&transaction)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		libs.HandleNotFoundError(c)
		return
	}

	c.JSON(http.StatusOK, models.PortfolioTransactionResponseFromDB(&transaction))
}
