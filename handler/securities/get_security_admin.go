package securities

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
	"gorm.io/gorm"
)

// GetSecurityAdmin returns security for user with admin privileges
func (h *securitiesHandler) GetSecurityAdmin(c *gin.Context) {
	uuid := c.Param("uuid")

	if err := h.Validate.Var(uuid, "required,uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var s db.Security
	err := h.DB.Preload("Events").Preload("SecurityMarkets").Preload("SecurityTaxonomies").Take(&s, "uuid = ?", uuid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			libs.HandleNotFoundError(c)
			return
		}

		panic(err)
	}

	events := []gin.H{}
	for _, e := range s.Events {
		events = append(events, gin.H{
			"id":           e.ID,
			"date":         e.Date,
			"type":         e.Type,
			"amount":       e.Amount,
			"currencyCode": e.CurrencyCode,
			"ratio":        e.Ratio,
			"securityUuid": e.SecurityUuid,
		})
	}

	markets := []gin.H{}
	for _, m := range s.SecurityMarkets {
		markets = append(markets, gin.H{
			"id":             m.ID,
			"securityUuid":   m.SecurityUUID,
			"marketCode":     m.MarketCode,
			"currencyCode":   m.CurrencyCode,
			"firstPriceDate": m.FirstPriceDate,
			"lastPriceDate":  m.LastPriceDate,
			"symbol":         m.Symbol,
			"updatePrices":   m.UpdatePrices,
		})
	}

	taxonomies := []gin.H{}
	for _, t := range s.SecurityTaxonomies {
		taxonomies = append(taxonomies, gin.H{
			"securityUuid": t.SecurityUUID,
			"taxonomyUuid": t.TaxonomyUUID,
			"weight":       t.Weight,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":               s.UUID,
		"name":               s.Name,
		"isin":               s.Isin,
		"wkn":                s.Wkn,
		"symbolXfra":         s.SymbolXfra,
		"symbolXnas":         s.SymbolXnas,
		"symbolXnys":         s.SymbolXnys,
		"securityType":       s.SecurityType,
		"markets":            markets,
		"events":             events,
		"securityTaxonomies": taxonomies,
	})
}
