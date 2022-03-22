package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
)

// DeleteSecurityMarket removes market of security
func (h *securitiesHandler) DeleteSecurityMarket(c *gin.Context) {

	uuid := c.Param("uuid")
	if err := h.validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}
	marketCode := c.Param("marketCode")

	market, err := h.SecurityService.DeleteSecurityMarket(uuid, marketCode)
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	c.JSON(http.StatusOK, market)
}
