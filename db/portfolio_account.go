package db

import "time"

// PortfolioAccount in database
type PortfolioAccount struct {
	PortfolioID          uint   `gorm:"primaryKey"`
	UUID                 string `gorm:"primaryKey"`
	Type                 string
	Name                 string
	CurrencyCode         *string
	ReferenceAccountUUID *string
	Active               bool
	Note                 string
	UpdatedAt            time.Time
}

// TableName defines name of table in database
func (PortfolioAccount) TableName() string {
	return "portfolios_accounts"
}
