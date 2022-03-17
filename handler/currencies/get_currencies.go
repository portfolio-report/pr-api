package currencies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/models"
)

// GetCurrencies returns all currencies with their exchange rates
func (h *currenciesHandler) GetCurrencies(c *gin.Context) {
	var currencies []db.Currency

	err := h.DB.
		Preload("ExchangeratesBase").
		Preload("ExchangeratesQuote").
		Find(&currencies).Error
	if err != nil {
		panic(err)
	}

	response := []models.CurrencyResponse{}
	for _, db := range currencies {
		response = append(response, models.CurrencyResponseFromDB(db))
	}

	c.JSON(http.StatusOK, response)
}
