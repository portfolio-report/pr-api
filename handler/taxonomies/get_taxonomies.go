package taxonomies

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Gets all taxonomies
func (h *TaxonomiesHandler) GetTaxonomies(c *gin.Context) {
	taxonomies, err := h.TaxonomyService.GetAllTaxonomies()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, taxonomies)
}
