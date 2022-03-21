package models

import (
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
)

type SecurityMarketResponsePublic struct {
	MarketCode     string      `json:"marketCode"`
	CurrencyCode   string      `json:"currencyCode"`
	FirstPriceDate *model.Date `json:"firstPriceDate"`
	LastPriceDate  *model.Date `json:"lastPriceDate"`
	Symbol         *string     `json:"symbol"`
}

func SecurityMarketResponsePublicFromDB(m *db.SecurityMarket) SecurityMarketResponsePublic {
	return SecurityMarketResponsePublic{
		MarketCode:     m.MarketCode,
		CurrencyCode:   m.CurrencyCode,
		Symbol:         m.Symbol,
		FirstPriceDate: m.FirstPriceDate,
		LastPriceDate:  m.LastPriceDate,
	}
}

type SecurityMarketPriceResponse struct {
	Date  model.Date       `json:"date"`
	Close db.DecimalString `json:"close"`
}

func SecurityMarketPriceResponseFromDB(p *db.SecurityMarketPrice) SecurityMarketPriceResponse {
	return SecurityMarketPriceResponse{
		Date:  p.Date,
		Close: p.Close,
	}
}
