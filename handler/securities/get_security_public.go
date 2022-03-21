package securities

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

// GetSecurityPublic returns security with its public attributes
func (h *securitiesHandler) GetSecurityPublic(c *gin.Context) {
	uuid := c.Param("uuid")

	if err := h.validate.Var(uuid, "required,LaxUuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var security db.Security
	err := h.DB.Take(&security, "uuid = ?", uuid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		libs.HandleNotFoundError(c)
		return
	}
	if err != nil {
		panic(err)
	}

	var markets []db.SecurityMarket
	err = h.DB.
		Select("market_code", "symbol", "currency_code", "first_price_date", "last_price_date").
		Where("security_uuid = ?", uuid).
		Find(&markets).Error
	if err != nil {
		panic(err)
	}

	var events []db.Event
	err = h.DB.
		Where("security_uuid = ? AND type IN ('dividend', 'split')", uuid).
		Find(&events).Error
	if err != nil {
		panic(err)
	}

	var securityTaxonomies []db.SecurityTaxonomy
	err = h.DB.
		Preload("Taxonomy").
		Where("security_uuid = ?", uuid).
		Find(&securityTaxonomies).Error
	if err != nil {
		panic(err)
	}

	eventsResp := []model.Event{}
	for _, e := range events {
		eventsResp = append(eventsResp, model.Event{
			Date:         e.Date.String(),
			Type:         e.Type,
			Amount:       e.Amount,
			Ratio:        e.Ratio,
			CurrencyCode: e.CurrencyCode,
		})
	}

	marketsResp := []models.SecurityMarketResponsePublic{}
	for _, m := range markets {
		marketsResp = append(marketsResp, models.SecurityMarketResponsePublicFromDB(&m))
	}

	taxonomiesResp := []models.SecurityTaxonomyResponse{}
	for _, t := range securityTaxonomies {
		taxonomiesResp = append(taxonomiesResp, models.SecurityTaxonomyResponseFromDB(&t))
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":               strings.Replace(security.UUID, "-", "", 4),
		"name":               security.Name,
		"isin":               security.Isin,
		"wkn":                security.Wkn,
		"symbolXfra":         security.SymbolXfra,
		"symbolXnas":         security.SymbolXnas,
		"symbolXnys":         security.SymbolXnys,
		"securityType":       security.SecurityType,
		"markets":            marketsResp,
		"events":             eventsResp,
		"securityTaxonomies": taxonomiesResp,
	})
}
