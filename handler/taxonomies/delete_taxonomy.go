package taxonomies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/libs"
)

// DeleteTaxonomy removes taxonomy
func (h *taxonomiesHandler) DeleteTaxonomy(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
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
