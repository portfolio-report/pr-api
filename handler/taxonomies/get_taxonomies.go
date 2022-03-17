package taxonomies

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTaxonomies lists all taxonomies
func (h *taxonomiesHandler) GetTaxonomies(c *gin.Context) {
	taxonomies, err := h.TaxonomyService.GetAllTaxonomies()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, taxonomies)
}
