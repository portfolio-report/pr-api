package db

import (
	"time"

	"gorm.io/datatypes"
)

// PortfolioSecurity in database
type PortfolioSecurity struct {
	PortfolioID   uint   `gorm:"primaryKey"`
	UUID          string `gorm:"primaryKey"`
	Name          string
	CurrencyCode  string
	Isin          string
	Wkn           string
	Symbol        string
	Active        bool
	Note          string
	SecurityUUID  *string
	UpdatedAt     time.Time
	Calendar      *string
	Feed          *string
	FeedUrl       *string
	LatestFeed    *string
	LatestFeedUrl *string
	Attributes    datatypes.JSON
	Events        datatypes.JSON
	Properties    datatypes.JSON
}

// TableName defines name of table in database
func (PortfolioSecurity) TableName() string {
	return "portfolios_securities"
}
