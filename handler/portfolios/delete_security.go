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

func (h *PortfoliosHandler) DeleteSecurity(c *gin.Context) {
	portfolioId := middleware.PortfolioFromContext(c).ID
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	// Delete transactions of security
	err := h.DB.
		Where("portfolio_id = ? AND portfolio_security_uuid = ?", portfolioId, uuid).
		Delete(&db.PortfolioTransaction{}).Error
	if err != nil {
		panic(err)
	}

	var security db.PortfolioSecurity
	result := h.DB.
		Clauses(clause.Returning{}).
		Where("portfolio_id = ? AND uuid = ?", portfolioId, uuid).
		Delete(&security)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		libs.HandleNotFoundError(c)
		return
	}

	c.JSON(http.StatusOK, models.PortfolioSecurityResponseFromDB(&security))
}
