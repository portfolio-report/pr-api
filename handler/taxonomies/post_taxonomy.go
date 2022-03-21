package taxonomies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
)

// PostTaxonomy creates taxonomy
func (h *taxonomiesHandler) PostTaxonomy(c *gin.Context) {
	var request model.TaxonomyInput
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	taxonomy, err := h.TaxonomyService.CreateTaxonomy(&request)
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, taxonomy)
}
