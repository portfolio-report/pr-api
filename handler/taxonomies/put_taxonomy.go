package taxonomies

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"gorm.io/gorm"
)

// PutTaxonomy updates taxonomy
func (h *taxonomiesHandler) PutTaxonomy(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var request model.TaxonomyInput
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	taxonomy, err := h.TaxonomyService.UpdateTaxonomy(uuid, &request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			libs.HandleNotFoundError(c)
			return
		}
		libs.HandleBadRequestError(c, err.Error())
	}

	c.JSON(http.StatusOK, taxonomy)
}
