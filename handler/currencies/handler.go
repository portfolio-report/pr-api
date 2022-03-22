package currencies

import (
	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
)

type currenciesHandler struct {
	model.CurrenciesService
}

// NewHandler creates new CurrenciesHandler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	UserService model.UserService,
	SessionService model.SessionService,
	CurrenciesService model.CurrenciesService,
) {
	h := &currenciesHandler{
		CurrenciesService: CurrenciesService,
	}

	g := R.Group("/currencies")

	g.GET("/", h.GetCurrencies)
	g.POST("/convert", h.Convert)
}
