package db

import "time"

type Clientupdate struct {
	ID        uint `gorm:"primaryKey"`
	Timestamp time.Time
	Version   string
	Country   *string
	Useragent *string
}

func (Clientupdate) TableName() string {
	return "clientupdates"
}
