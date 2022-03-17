package taxonomies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
)

// GetTaxonomy returns single taxonomy with all descendants
func (h *taxonomiesHandler) GetTaxonomy(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	taxonomy, err := h.TaxonomyService.GetTaxonomyByUUID(uuid)
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	descendants, err := h.TaxonomyService.GetDescendantsOfTaxonomy(taxonomy)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":        taxonomy.UUID,
		"parentUuid":  taxonomy.ParentUUID,
		"rootUuid":    taxonomy.RootUUID,
		"name":        taxonomy.Name,
		"code":        taxonomy.Code,
		"descendants": descendants,
	})
}
