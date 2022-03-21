package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/shopspring/decimal"
)

type securityTaxonomyRequest struct {
	TaxonomyUuid string          `json:"taxonomyUuid" binding:"uuid"`
	Weight       decimal.Decimal `json:"weight"`
}

// PutSecurityTaxonomies creates, updates and deletes taxonomies of security
func (h *securitiesHandler) PutSecurityTaxonomies(c *gin.Context) {
	securityUuid := c.Param("uuid")
	if err := h.validate.Var(securityUuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}
	taxonomyRootUuid := c.Param("rootUuid")
	if err := h.validate.Var(taxonomyRootUuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var req []securityTaxonomyRequest
	if err := c.BindJSON(&req); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	// Convert []securityTaxonomyRequest to []*model.SecurityTaxonomyInput
	inputs := make([]*model.SecurityTaxonomyInput, len(req))
	for i := range req {
		inputs[i] = &model.SecurityTaxonomyInput{
			TaxonomyUUID: req[i].TaxonomyUuid,
			Weight:       req[i].Weight,
		}
	}

	ret, err := h.SecurityService.UpdateSecurityTaxonomies(securityUuid, taxonomyRootUuid, inputs)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ret)
}
