package db

import "time"

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

func (PortfolioAccount) TableName() string {
	return "portfolios_accounts"
}
