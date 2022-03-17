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

// DeleteAccount removes account from portfolio and links to it
func (h *PortfoliosHandler) DeleteAccount(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	// Remove links as reference account
	h.DB.Model(&db.PortfolioAccount{}).
		Where("portfolio_id = ? AND reference_account_uuid = ?", portfolioId, uuid).
		Update("reference_account_uuid", nil)

	// Delete transactions of account
	err := h.DB.
		Where("portfolio_id = ? AND account_uuid = ?", portfolioId, uuid).
		Delete(&db.PortfolioTransaction{}).Error
	if err != nil {
		panic(err)
	}

	var account db.PortfolioAccount
	result := h.DB.Clauses(clause.Returning{}).Where("portfolio_id = ? AND uuid = ?", portfolioId, uuid).Delete(&account)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		libs.HandleNotFoundError(c)
		return
	}

	c.JSON(http.StatusOK, models.PortfolioAccountResponseFromDB(&account))
}
