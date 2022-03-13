package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm/clause"
)

// Deletes security
func (h *SecuritiesHandler) DeleteSecurity(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var security db.Security
	result := h.DB.Clauses(clause.Returning{}).Delete(&security, "uuid = ?", uuid)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		libs.HandleNotFoundError(c)
		return
	}

	c.JSON(http.StatusOK, models.SecurityResponseFromDB(&security))
}
