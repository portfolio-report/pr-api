package taxonomies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
)

type postTaxonomyRequest struct {
	Name       string  `json:"name" binding:"required"`
	Code       *string `json:"code"`
	ParentUuid *string `json:"parentUuid" binding:"omitempty,uuid"`
}

// Creates taxonomy
func (h *TaxonomiesHandler) PostTaxonomy(c *gin.Context) {
	var request postTaxonomyRequest
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	taxonomy, err := h.TaxonomyService.CreateTaxonomy(&model.Taxonomy{
		Name:       request.Name,
		Code:       request.Code,
		ParentUUID: request.ParentUuid,
	})
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, taxonomy)
}
