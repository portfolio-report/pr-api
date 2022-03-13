package db

import (
	"time"
)

type Portfolio struct {
	ID               uint `gorm:"primaryKey"`
	Name             string
	Note             string
	BaseCurrencyCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UserID           uint
}

func (Portfolio) TableName() string {
	return "portfolios"
}
