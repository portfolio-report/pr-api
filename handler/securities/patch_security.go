package securities

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

// PatchSecurity updates security
func (h *SecuritiesHandler) PatchSecurity(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var request CreateUpdateSecurityRequest
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	var security db.Security
	if err := h.DB.Take(&security, "uuid = ?", uuid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			libs.HandleNotFoundError(c)
			return
		}
		panic(err)
	}

	if request.Name != nil {
		security.Name = request.Name
	}
	if request.Isin != nil {
		security.Isin = request.Isin
	}
	if request.Wkn != nil {
		security.Wkn = request.Wkn
	}
	if request.SecurityType != nil {
		security.SecurityType = request.SecurityType
	}
	if request.SymbolXfra != nil {
		security.SymbolXfra = request.SymbolXfra
	}
	if request.SymbolXnas != nil {
		security.SymbolXnas = request.SymbolXnas
	}
	if request.SymbolXnys != nil {
		security.SymbolXnys = request.SymbolXnys
	}

	if err := h.DB.Save(&security).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, models.SecurityResponseFromDB(&security))
}
