package db

type Currency struct {
	Code               string         `gorm:"primaryKey"`
	ExchangeratesBase  []Exchangerate `gorm:"foreignKey:base_currency_code;references:code"`
	ExchangeratesQuote []Exchangerate `gorm:"foreignKey:quote_currency_code;references:code"`
}

func (Currency) TableName() string {
	return "currencies"
}
