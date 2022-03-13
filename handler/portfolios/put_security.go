package portfolios

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
)

type putSecurityRequest struct {
	Name          string    `json:"name" binding:"required"`
	CurrencyCode  string    `json:"currencyCode" binding:"max=3"`
	Isin          string    `json:"isin"`
	Wkn           string    `json:"wkn"`
	Symbol        string    `json:"symbol"`
	Active        bool      `json:"active"`
	Note          string    `json:"note"`
	SecurityUuid  *string   `json:"securityUuid" binding:"omitempty,LaxUuid"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Calendar      *string   `json:"calendar"`
	Feed          *string   `json:"feed"`
	FeedUrl       *string   `json:"feedUrl"`
	LatestFeed    *string   `json:"latestFeed"`
	LatestFeedUrl *string   `json:"latestFeedUrl"`

	Events []struct {
		Date    string `json:"date" binding:"DateYYYY-MM-DD"`
		Type    string `json:"type" binding:"oneof=STOCK_SPLIT NOTE DIVIDEND_PAYMENT"`
		Details string `json:"details"`
	} `json:"events" binding:"dive"`
	Properties []struct {
		Name  string `json:"name"`
		Type  string `json:"type" binding:"oneof=MARKET FEED"`
		Value string `json:"value"`
	} `json:"properties" binding:"dive"`
}

func (h *PortfoliosHandler) PutSecurity(c *gin.Context) {
	portfolioId := uint(middleware.PortfolioFromContext(c).ID)
	uuid := c.Param("uuid")
	if err := h.Validate.Var(uuid, "uuid"); err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var req putSecurityRequest
	if err := c.BindJSON(&req); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	var security db.PortfolioSecurity

	err := h.DB.FirstOrInit(&security, db.PortfolioSecurity{PortfolioID: portfolioId, UUID: uuid}).Error
	if err != nil {
		panic(err)
	}

	security.Name = req.Name
	security.CurrencyCode = req.CurrencyCode
	security.Isin = req.Isin
	security.Wkn = req.Wkn
	security.Symbol = req.Symbol
	security.Active = req.Active
	security.Note = req.Note
	security.SecurityUUID = req.SecurityUuid
	security.UpdatedAt = req.UpdatedAt
	security.Calendar = req.Calendar
	security.Feed = req.Feed
	security.FeedUrl = req.FeedUrl
	security.LatestFeed = req.LatestFeed
	security.LatestFeedUrl = req.LatestFeedUrl
	security.Events, err = json.Marshal(req.Events)
	if err != nil {
		panic(err)
	}
	security.Properties, err = json.Marshal(req.Properties)
	if err != nil {
		panic(err)
	}

	if err := h.DB.Save(&security).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			libs.HandleBadRequestError(c, "data violates constraint "+pqErr.Constraint)
			return
		}

		panic(err)
	}

	c.JSON(http.StatusOK, models.PortfolioSecurityResponseFromDB(&security))
}
