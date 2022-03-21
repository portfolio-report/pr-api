package securities

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"gorm.io/gorm/clause"
)

type patchSecurityMarketRequest struct {
	CurrencyCode *string `json:"currencyCode"`
	Symbol       *string `json:"symbol"`
	UpdatePrices *bool   `json:"updatePrices"`
	Prices       *[]struct {
		Date  model.Date `json:"date"`
		Close float64    `json:"close"`
	} `json:"prices"`
}

// PatchSecurityMarket creates or updates market of security and its prices
func (h *securitiesHandler) PatchSecurityMarket(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}
	marketCode := c.Param("marketCode")

	var req patchSecurityMarketRequest
	if err := c.BindJSON(&req); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	var market db.SecurityMarket
	err := h.DB.
		Attrs(db.SecurityMarket{UpdatePrices: true}).
		FirstOrInit(&market, db.SecurityMarket{SecurityUUID: uuid, MarketCode: marketCode}).
		Error
	if err != nil {
		panic(err)
	}

	if req.CurrencyCode != nil {
		market.CurrencyCode = *req.CurrencyCode
	}
	if req.Symbol != nil {
		market.Symbol = req.Symbol
	}
	if req.UpdatePrices != nil {
		market.UpdatePrices = *req.UpdatePrices
	}

	if err := h.DB.Save(&market).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			libs.HandleBadRequestError(c, "data violates constraint "+pqErr.Constraint)
			return
		}

		panic(err)
	}

	// Create/update the associated prices
	if req.Prices != nil {
		prices := []db.SecurityMarketPrice{}
		for _, p := range *req.Prices {
			prices = append(prices, db.SecurityMarketPrice{
				SecurityMarketID: market.ID,
				Date:             p.Date,
				Close:            db.DecimalString(strconv.FormatFloat(p.Close, 'f', 8, 64)),
			})
		}

		if err := h.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&prices).Error; err != nil {
			panic(err)
		}
	}

	// Keep firstPriceDate and lastPriceDate up-to-date
	err = h.DB.Exec(`UPDATE securities_markets SET `+
		`first_price_date = (SELECT MIN(date) FROM securities_markets_prices WHERE security_market_id = ?), `+
		`last_price_date =  (SELECT MAX(date) FROM securities_markets_prices WHERE security_market_id = ?) `+
		`WHERE id = ?`, market.ID, market.ID, market.ID).Error
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
