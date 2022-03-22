package securities

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"gorm.io/gorm"
)

// PatchSecurity updates security
func (h *securitiesHandler) PatchSecurity(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var request model.SecurityInput
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	security, err := h.SecurityService.GetSecurityByUUID(uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			libs.HandleNotFoundError(c)
			return
		}
		panic(err)
	}

	if request.Name == nil {
		request.Name = security.Name
	}
	if request.Isin == nil {
		request.Isin = security.Isin
	}
	if request.Wkn == nil {
		request.Wkn = security.Wkn
	}
	if request.SecurityType == nil {
		request.SecurityType = security.SecurityType
	}
	if request.SymbolXfra == nil {
		request.SymbolXfra = security.SymbolXfra
	}
	if request.SymbolXnas == nil {
		request.SymbolXnas = security.SymbolXnas
	}
	if request.SymbolXnys == nil {
		request.SymbolXnys = security.SymbolXnys
	}

	security, err = h.SecurityService.UpdateSecurity(uuid, &request)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, security)
}
