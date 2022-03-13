package models

import (
	"time"

	"github.com/portfolio-report/pr-api/db"
	"gorm.io/datatypes"
)

type PortfolioSecurityResponse struct {
	UUID          string         `json:"uuid"`
	Name          string         `json:"name"`
	CurrencyCode  string         `json:"currencyCode"`
	Isin          string         `json:"isin"`
	Wkn           string         `json:"wkn"`
	Symbol        string         `json:"symbol"`
	Active        bool           `json:"active"`
	Note          string         `json:"note"`
	SecurityUUID  *string        `json:"securityUuid"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	Calendar      *string        `json:"calendar"`
	Feed          *string        `json:"feed"`
	FeedUrl       *string        `json:"feedUrl"`
	LatestFeed    *string        `json:"latestFeed"`
	LatestFeedUrl *string        `json:"latestFeedUrl"`
	Events        datatypes.JSON `json:"events"`
	Properties    datatypes.JSON `json:"properties"`
}

func PortfolioSecurityResponseFromDB(s *db.PortfolioSecurity) PortfolioSecurityResponse {
	return PortfolioSecurityResponse{
		UUID:          s.UUID,
		Name:          s.Name,
		CurrencyCode:  s.CurrencyCode,
		Isin:          s.Isin,
		Wkn:           s.Wkn,
		Symbol:        s.Symbol,
		Active:        s.Active,
		Note:          s.Note,
		SecurityUUID:  s.SecurityUUID,
		UpdatedAt:     s.UpdatedAt.UTC(),
		Calendar:      s.Calendar,
		Feed:          s.Feed,
		FeedUrl:       s.FeedUrl,
		LatestFeed:    s.LatestFeed,
		LatestFeedUrl: s.LatestFeedUrl,
		Events:        s.Events,
		Properties:    s.Properties,
	}
}
