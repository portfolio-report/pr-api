package taxonomies

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTaxonomies lists all taxonomies
func (h *taxonomiesHandler) GetTaxonomies(c *gin.Context) {
	taxonomies := h.TaxonomyService.GetAllTaxonomies()
	c.JSON(http.StatusOK, taxonomies)
}
