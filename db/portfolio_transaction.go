package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PortfolioTransaction in database
type PortfolioTransaction struct {
	PortfolioID            uint      `gorm:"primaryKey"`
	UUID                   uuid.UUID `gorm:"primaryKey"`
	AccountUUID            uuid.UUID
	Type                   string
	Datetime               time.Time
	PartnerTransactionUUID *uuid.UUID
	Shares                 decimal.NullDecimal
	PortfolioSecurityUUID  *uuid.UUID
	Note                   string
	UpdatedAt              time.Time

	Units []PortfolioTransactionUnit `gorm:"foreignKey:portfolio_id,transaction_uuid;references:portfolio_id,uuid"`
}

// TableName defines name of table in database
func (PortfolioTransaction) TableName() string {
	return "portfolios_transactions"
}

// PortfolioTransactionUnit in database
type PortfolioTransactionUnit struct {
	ID                   uint `gorm:"primaryKey"`
	TransactionUUID      uuid.UUID
	PortfolioID          uint
	Type                 string
	Amount               decimal.Decimal
	CurrencyCode         string
	OriginalAmount       decimal.NullDecimal
	OriginalCurrencyCode *string
	ExchangeRate         decimal.NullDecimal
}

// TableName defines name of table in database
func (PortfolioTransactionUnit) TableName() string {
	return "portfolios_transactions_units"
}
