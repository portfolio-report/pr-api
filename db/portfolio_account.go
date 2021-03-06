package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/graph/model"
)

// PortfolioAccount in database
type PortfolioAccount struct {
	PortfolioID          uint      `gorm:"primaryKey"`
	UUID                 uuid.UUID `gorm:"primaryKey"`
	Type                 model.PortfolioAccountType
	Name                 string
	CurrencyCode         *string
	ReferenceAccountUUID *uuid.UUID
	Active               bool
	Note                 string
	UpdatedAt            time.Time
}

// TableName defines name of table in database
func (PortfolioAccount) TableName() string {
	return "portfolios_accounts"
}
