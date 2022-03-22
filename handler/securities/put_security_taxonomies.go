package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/shopspring/decimal"
)

type securityTaxonomyRequest struct {
	TaxonomyUuid uuid.UUID       `json:"taxonomyUuid"`
	Weight       decimal.Decimal `json:"weight"`
}

// PutSecurityTaxonomies creates, updates and deletes taxonomies of security
func (h *securitiesHandler) PutSecurityTaxonomies(c *gin.Context) {
	securityUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}
	taxonomyRootUuid, err := uuid.Parse(c.Param("rootUuid"))
	if err != nil {
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
