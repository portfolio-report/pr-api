package db

type Exchangerate struct {
	ID                uint `gorm:"primaryKey"`
	BaseCurrencyCode  string
	QuoteCurrencyCode string
}

func (Exchangerate) TableName() string {
	return "exchangerates"
}

type ExchangeratePrice struct {
	ExchangerateID uint   `gorm:"primaryKey;autoIncrement:false"`
	Date           DbDate `gorm:"primaryKey"`
	Value          DecimalString
}

func (ExchangeratePrice) TableName() string {
	return "exchangerates_prices"
}
