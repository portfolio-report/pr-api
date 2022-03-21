package db

import "github.com/portfolio-report/pr-api/graph/model"

// Exchangerate in database
type Exchangerate struct {
	ID                uint `gorm:"primaryKey"`
	BaseCurrencyCode  string
	QuoteCurrencyCode string
}

// TableName defines name of table in database
func (Exchangerate) TableName() string {
	return "exchangerates"
}

// ExchangeratePrice in database
type ExchangeratePrice struct {
	ExchangerateID uint       `gorm:"primaryKey;autoIncrement:false"`
	Date           model.Date `gorm:"primaryKey"`
	Value          DecimalString
}

// TableName defines name of table in database
func (ExchangeratePrice) TableName() string {
	return "exchangerates_prices"
}
