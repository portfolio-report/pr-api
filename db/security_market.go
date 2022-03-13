package db

type SecurityMarket struct {
	ID             uint `gorm:"primaryKey"`
	SecurityUUID   string
	MarketCode     string
	CurrencyCode   string
	FirstPriceDate *DbDate
	LastPriceDate  *DbDate
	Symbol         *string
	UpdatePrices   bool
}

func (SecurityMarket) TableName() string {
	return "securities_markets"
}

type SecurityMarketPrice struct {
	SecurityMarketID uint   `gorm:"primaryKey;autoIncrement:false"`
	Date             DbDate `gorm:"primaryKey"`
	Close            DecimalString
}

func (SecurityMarketPrice) TableName() string {
	return "securities_markets_prices"
}
