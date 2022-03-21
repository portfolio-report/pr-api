package currencies

import (
	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/models"
)

type currenciesHandler struct {
	CurrenciesService models.CurrenciesService
}

// NewHandler creates new CurrenciesHandler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	UserService models.UserService,
	SessionService models.SessionService,
	CurrenciesService models.CurrenciesService,
) {
	h := &currenciesHandler{
		CurrenciesService: CurrenciesService,
	}

	g := R.Group("/currencies")

	g.GET("/", h.GetCurrencies)
	g.POST("/convert", h.Convert)
}
