package currencies

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCurrencies returns all currencies with their exchange rates
func (h *currenciesHandler) GetCurrencies(c *gin.Context) {
	currencies := h.CurrenciesService.GetCurrencies()
	c.JSON(http.StatusOK, currencies)
}
