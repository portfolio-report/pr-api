package db

import "github.com/portfolio-report/pr-api/graph/model"

// SecurityMarket in database
type SecurityMarket struct {
	ID             uint `gorm:"primaryKey"`
	SecurityUUID   string
	MarketCode     string
	CurrencyCode   string
	FirstPriceDate *model.Date
	LastPriceDate  *model.Date
	Symbol         *string
	UpdatePrices   bool
}

// TableName defines name of table in database
func (SecurityMarket) TableName() string {
	return "securities_markets"
}

// SecurityMarketPrice in database
type SecurityMarketPrice struct {
	SecurityMarketID uint       `gorm:"primaryKey;autoIncrement:false"`
	Date             model.Date `gorm:"primaryKey"`
	Close            DecimalString
}

// TableName defines name of table in database
func (SecurityMarketPrice) TableName() string {
	return "securities_markets_prices"
}
