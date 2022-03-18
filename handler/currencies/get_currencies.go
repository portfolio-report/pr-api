package currencies

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCurrencies returns all currencies with their exchange rates
func (h *currenciesHandler) GetCurrencies(c *gin.Context) {
	currencies, err := h.CurrenciesService.GetCurrencies()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, currencies)
}
