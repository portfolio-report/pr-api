package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm/clause"
)

// DeleteSecurityMarket removes market of security
func (h *SecuritiesHandler) DeleteSecurityMarket(c *gin.Context) {

	uuid := c.Param("uuid")
	if err := h.validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}
	marketCode := c.Param("marketCode")

	var market db.SecurityMarket
	result := h.DB.Clauses(clause.Returning{}).Delete(&market, "security_uuid = ? AND market_code = ?", uuid, marketCode)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		libs.HandleNotFoundError(c)
		return
	}

	c.JSON(http.StatusOK, models.SecurityMarketResponsePublicFromDB(&market))
}
