package currencies

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/shopspring/decimal"
)

type convertRequest struct {
	SourceCurrencyCode string          `json:"sourceCurrencyCode" binding:"required"`
	TargetCurrencyCode string          `json:"targetCurrencyCode" binding:"required"`
	SourceAmount       decimal.Decimal `json:"sourceAmount"`
	Date               *model.Date     `json:"date"`
}

// Convert converts amount between currencies
func (h *currenciesHandler) Convert(c *gin.Context) {
	var r convertRequest
	if err := c.BindJSON(&r); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	if r.Date == nil {
		now := model.Date(time.Now())
		r.Date = &now
	}

	targetAmount, err := h.CurrenciesService.ConvertCurrencyAmount(r.SourceAmount, r.SourceCurrencyCode, r.TargetCurrencyCode, r.Date.Time())
	if err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sourceAmount":       r.SourceAmount,
		"targetAmount":       targetAmount,
		"sourceCurrencyCode": r.SourceCurrencyCode,
		"targetCurrencyCode": r.TargetCurrencyCode,
		"date":               r.Date})
}
