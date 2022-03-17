package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
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

	// Remove securityTaxonomies not in request
	taxonomyUuids := []string{}
	for _, st := range req {
		taxonomyUuids = append(taxonomyUuids, st.TaxonomyUuid)
	}
	err := h.DB.Exec("DELETE FROM securities_taxonomies st "+
		"USING taxonomies t "+
		"WHERE st.taxonomy_uuid = t.uuid"+
		" AND st.security_uuid = ?"+
		" AND t.root_uuid = ?"+
		" AND st.taxonomy_uuid NOT IN ?", securityUuid, taxonomyRootUuid, taxonomyUuids).
		Error
	if err != nil {
		panic(err)
	}

	// Upsert all securityTaxonomies in request
	securityTaxonomies := []db.SecurityTaxonomy{}
	for _, r := range req {
		securityTaxonomies = append(securityTaxonomies, db.SecurityTaxonomy{
			SecurityUUID: securityUuid,
			TaxonomyUUID: r.TaxonomyUuid,
			Weight:       r.Weight,
		})
	}
	if err := h.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&securityTaxonomies).Error; err != nil {
		panic(err)
	}

	ret := []models.SecurityTaxonomyResponse{}
	for _, t := range securityTaxonomies {
		ret = append(ret, models.SecurityTaxonomyResponseFromDB(&t))
	}
	c.JSON(http.StatusOK, ret)
}
