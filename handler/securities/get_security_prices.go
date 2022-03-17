package securities

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

// GetSecurityPrices returns the market and prices of security
func (h *SecuritiesHandler) GetSecurityPrices(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.validate.Var(uuid, "LaxUuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	from := c.Query("from")
	if from == "" {
		from = time.Now().AddDate(0, 0, -14).Format("2006-01-02")
	}

	if err := h.validate.Var(from, "DateYYYY-MM-DD"); err != nil {
		libs.HandleBadRequestError(c, "from is not a valid date")
		return
	}

	var market db.SecurityMarket
	var prices []db.SecurityMarketPrice

	err := h.DB.
		Where("market_code = 'XETR' AND security_uuid = ?", uuid).
		Take(&market).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		libs.HandleNotFoundError(c)
		return
	}
	if err != nil {
		panic(err)
	}

	err = h.DB.
		Select("date", "close").Where("security_market_id = ? AND date >= ?", market.ID, from).
		Order("date").
		Find(&prices).Error

	if err != nil {
		panic(err)
	}

	pricesResponse := []models.SecurityMarketPriceResponse{}
	for _, p := range prices {
		pricesResponse = append(pricesResponse, models.SecurityMarketPriceResponseFromDB(&p))
	}

	c.JSON(http.StatusOK, gin.H{
		"marketCode":     market.MarketCode,
		"currencyCode":   market.CurrencyCode,
		"symbol":         market.Symbol,
		"firstPriceDate": market.FirstPriceDate,
		"lastPriceDate":  market.LastPriceDate,
		"prices":         pricesResponse})

}
