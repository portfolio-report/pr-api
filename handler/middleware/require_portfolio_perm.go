package middleware

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/models"

	"github.com/portfolio-report/pr-api/libs"
	"gorm.io/gorm"
)

// RequirePortfolioPerm returns middleware which checks if URL parameter porfolioId
// belongs to the current user and stores the portfolio in Gin context
func RequirePortfolioPerm(PortfolioService models.PortfolioService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserFromContext(c.Request.Context())
		portfolioId, err := strconv.Atoi(c.Param("portfolioId"))
		if err != nil {
			libs.HandleNotFoundError(c)
			return
		}

		portfolio, err := PortfolioService.GetPortfolioOfUserByID(user, uint(portfolioId))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				libs.HandleNotFoundError(c)
				return
			}

			panic(err)
		}

		c.Set("portfolio", portfolio)

		c.Next()
	}
}

// PortfolioFromContext gets portfolio from Gin context
func PortfolioFromContext(c *gin.Context) *model.Portfolio {
	return c.MustGet("portfolio").(*model.Portfolio)
}
