package models

import (
	"github.com/portfolio-report/pr-api/db"
	"github.com/shopspring/decimal"
)

type SecurityTaxonomyResponse struct {
	TaxonomyUUID     string
	Weight           decimal.Decimal
	RootTaxonomyUUID *string `json:"rootTaxonomyUuid"`
}

func SecurityTaxonomyResponseFromDB(st *db.SecurityTaxonomy) SecurityTaxonomyResponse {
	return SecurityTaxonomyResponse{
		TaxonomyUUID:     st.TaxonomyUUID,
		Weight:           st.Weight,
		RootTaxonomyUUID: st.Taxonomy.RootUUID,
	}
}

type SecurityMarketResponsePublic struct {
	MarketCode     string     `json:"marketCode"`
	CurrencyCode   string     `json:"currencyCode"`
	FirstPriceDate *db.DbDate `json:"firstPriceDate"`
	LastPriceDate  *db.DbDate `json:"lastPriceDate"`
	Symbol         *string    `json:"symbol"`
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
	Date  db.DbDate        `json:"date"`
	Close db.DecimalString `json:"close"`
}

func SecurityMarketPriceResponseFromDB(p *db.SecurityMarketPrice) SecurityMarketPriceResponse {
	return SecurityMarketPriceResponse{
		Date:  p.Date,
		Close: p.Close,
	}
}
