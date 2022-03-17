package db

import "time"

// Clientupdate in database
type Clientupdate struct {
	ID        uint `gorm:"primaryKey"`
	Timestamp time.Time
	Version   string
	Country   *string
	Useragent *string
}

// TableName defines name of table in database
func (Clientupdate) TableName() string {
	return "clientupdates"
}
