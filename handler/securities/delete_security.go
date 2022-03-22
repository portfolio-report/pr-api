package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
)

// DeleteSecurity removes security
func (h *securitiesHandler) DeleteSecurity(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	security, err := h.SecurityService.DeleteSecurity(uuid)
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	c.JSON(http.StatusOK, security)
}
