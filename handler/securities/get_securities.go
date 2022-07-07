package securities

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
)

// GetSecurities lists securities
func (h *securitiesHandler) GetSecurities(c *gin.Context) {
	type Query struct {
		Limit        int    `form:"limit"`
		Skip         int    `form:"skip"`
		Sort         string `form:"sort" binding:"omitempty,oneof=uuid name isin wkn symbolXfra symbolXnas symbolXnys"`
		Desc         bool   `form:"desc"`
		Search       string `form:"search"`
		SecurityType string `form:"securityType"`
	}

	var q Query

	if err := c.BindQuery(&q); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	if q.Limit == 0 {
		q.Limit = 10
	}

	switch q.Sort {
	case "symbolXfra":
		q.Sort = "symbol_xfra"
	case "symbolXnas":
		q.Sort = "symbol_xnas"
	case "symbolXnys":
		q.Sort = "symbol_xnys"
	}

	query := h.DB

	if q.Search != "" {
		if err := h.Validate.Var(q.Search, "required,LaxUuid"); err == nil {
			query = query.Where("uuid = ?", q.Search)
		} else {
			like := "%" + q.Search + "%"
			query = query.Where(
				"name ILIKE ? OR isin ILIKE ? OR wkn ILIKE ? OR "+
					"symbol_xfra ILIKE ? OR symbol_xnas ILIKE ? OR symbol_xnys ILIKE ?",
				like, like, like, like, like, like)
		}
	}

	if q.SecurityType != "" {
		query = query.Where("security_type = ?", q.SecurityType)
	}

	var totalCount int64
	if err := query.Table("securities").Count(&totalCount).Error; err != nil {
		panic(err)
	}

	order := "name"
	for _, col := range []string{"uuid", "name", "isin", "wkn", "symbol_xfra", "symbol_xnas", "symbol_xnys"} {
		if q.Sort == col {
			order = col
		}
	}
	if q.Desc {
		order += " desc"
	}
	query = query.Order(order)

	query = query.Limit(q.Limit).Offset(q.Skip)

	var securities []db.Security
	if err := query.Preload("Events").Preload("SecurityMarkets").Find(&securities).Error; err != nil {
		panic(err)
	}

	entries := []gin.H{}
	for _, s := range securities {
		markets := []gin.H{}
		for _, m := range s.SecurityMarkets {
			markets = append(markets, gin.H{
				"marketCode":     m.MarketCode,
				"currencyCode":   m.CurrencyCode,
				"firstPriceDate": m.FirstPriceDate,
				"lastPriceDate":  m.LastPriceDate,
				"symbol":         m.Symbol,
				"updatePrices":   m.UpdatePrices,
			})
		}

		events := []gin.H{}
		for _, e := range s.Events {
			if e.Type == "dividend" || e.Type == "split" {
				events = append(events, gin.H{
					"date":         e.Date,
					"type":         e.Type,
					"amount":       e.Amount,
					"currencyCode": e.CurrencyCode,
					"ratio":        e.Ratio,
				})
			}
		}

		logoUrl := h.SecurityService.LogoUrlFromExtras(s.Extras)

		entries = append(entries, gin.H{
			"uuid":         s.UUID,
			"name":         s.Name,
			"isin":         s.Isin,
			"wkn":          s.Wkn,
			"symbolXfra":   s.SymbolXfra,
			"symbolXnas":   s.SymbolXnas,
			"symbolXnys":   s.SymbolXnys,
			"securityType": s.SecurityType,
			"markets":      markets,
			"events":       events,
			"logoUrl":      logoUrl,
		})
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries, "params": gin.H{"totalCount": totalCount}})
}
