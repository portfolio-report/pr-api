package taxonomies

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"gorm.io/gorm"
)

type patchTaxonomyRequest struct {
	Name       *string `json:"name"`
	Code       *string `json:"code"`
	ParentUuid *string `json:"parentUuid" binding:"omitempty,uuid"`
}

// PatchTaxonomy updates taxonomy
func (h *taxonomiesHandler) PatchTaxonomy(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var request patchTaxonomyRequest
	if err := c.BindJSON(&request); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	taxonomy := &model.Taxonomy{
		Name:       *request.Name,
		Code:       request.Code,
		ParentUUID: request.ParentUuid,
	}

	taxonomy, err := h.TaxonomyService.UpdateTaxonomy(uuid, taxonomy)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			libs.HandleNotFoundError(c)
			return
		}
		libs.HandleBadRequestError(c, err.Error())
	}

	c.JSON(http.StatusOK, taxonomy)
}
