package db

import "github.com/portfolio-report/pr-api/graph/model"

// Event in database
type Event struct {
	ID           uint `gorm:"primaryKey"`
	Date         model.Date
	Type         string
	Amount       *string
	CurrencyCode *string
	Ratio        *string
	SecurityUuid string
}

// TableName defines name of table in database
func (Event) TableName() string {
	return "events"
}
