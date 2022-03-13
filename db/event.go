package db

type Event struct {
	ID           uint `gorm:"primaryKey"`
	Date         DbDate
	Type         string
	Amount       *string
	CurrencyCode *string
	Ratio        *string
	SecurityUuid string
}

func (Event) TableName() string {
	return "events"
}
