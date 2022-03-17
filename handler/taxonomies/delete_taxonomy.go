package taxonomies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
)

// DeleteTaxonomy removes taxonomy
func (h *TaxonomiesHandler) DeleteTaxonomy(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	taxonomy, err := h.TaxonomyService.DeleteTaxonomy(uuid)

	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	c.JSON(http.StatusOK, taxonomy)
}
