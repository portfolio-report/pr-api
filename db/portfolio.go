package db

import (
	"time"
)

// Portfolio in database
type Portfolio struct {
	ID               uint `gorm:"primaryKey"`
	Name             string
	Note             string
	BaseCurrencyCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UserID           uint
}

// TableName defines name of table in database
func (Portfolio) TableName() string {
	return "portfolios"
}
