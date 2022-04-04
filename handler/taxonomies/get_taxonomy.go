package taxonomies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/libs"
)

// GetTaxonomy returns single taxonomy with all descendants
func (h *taxonomiesHandler) GetTaxonomy(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	taxonomy, err := h.TaxonomyService.GetTaxonomyByUUID(uuid)
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	descendants := h.TaxonomyService.GetDescendantsOfTaxonomy(taxonomy)

	c.JSON(http.StatusOK, gin.H{
		"uuid":        taxonomy.UUID,
		"parentUuid":  taxonomy.ParentUUID,
		"rootUuid":    taxonomy.RootUUID,
		"name":        taxonomy.Name,
		"code":        taxonomy.Code,
		"descendants": descendants,
	})
}
