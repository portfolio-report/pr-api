package currencies

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

type getExchangerateQuery struct {
	StartDate string `form:"startDate" binding:"omitempty,DateYYYY-MM-DD"`
}

// GetExchangerate returns a single exchange rate and its prices
func (h *CurrenciesHandler) GetExchangerate(c *gin.Context) {
	baseCurrencyCode := c.Param("baseCurrencyCode")
	quoteCurrencyCode := c.Param("quoteCurrencyCode")

	var q getExchangerateQuery
	if err := c.BindQuery(&q); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	if q.StartDate == "" {
		q.StartDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	var exchangerate db.Exchangerate
	err := h.DB.
		Where("base_currency_code = ? AND quote_currency_code = ?", baseCurrencyCode, quoteCurrencyCode).
		Take(&exchangerate).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		libs.HandleNotFoundError(c)
		return
	}
	if err != nil {
		panic(err)
	}

	var prices []db.ExchangeratePrice
	err = h.DB.
		Where("date >= ?", q.StartDate).
		Where("exchangerate_id = ?", exchangerate.ID).Order("date ASC").Find(&prices).Error
	if err != nil {
		panic(err)
	}

	pricesResponse := []models.ExchangeratePriceResponse{}
	for _, p := range prices {
		pricesResponse = append(pricesResponse, models.ExchangeratePriceResponseFromDB(p))
	}

	var latestPriceDate time.Time
	err = h.DB.
		Raw("SELECT max(date) FROM exchangerates_prices WHERE exchangerate_id = ?", exchangerate.ID).
		Scan(&latestPriceDate).Error
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"baseCurrencyCode":  exchangerate.BaseCurrencyCode,
		"quoteCurrencyCode": exchangerate.QuoteCurrencyCode,
		"latestPriceDate":   latestPriceDate.Format("2006-01-02"),
		"prices":            pricesResponse,
	})
}
