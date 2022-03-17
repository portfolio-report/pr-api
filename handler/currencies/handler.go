package currencies

import (
	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
)

type CurrenciesHandler struct {
	DB                *gorm.DB
	CurrenciesService models.CurrenciesService
}

// NewHandler creates new CurrenciesHandler and registers routes
func NewHandler(
	R *gin.RouterGroup,
	DB *gorm.DB,
	UserService models.UserService,
	SessionService models.SessionService,
	CurrenciesService models.CurrenciesService,
) {
	h := &CurrenciesHandler{
		DB:                DB,
		CurrenciesService: CurrenciesService,
	}

	g := R.Group("/currencies")

	g.GET("/", h.GetCurrencies)
	g.GET("/:baseCurrencyCode/:quoteCurrencyCode", h.GetExchangerate)
	g.POST("/convert", h.Convert)
}
