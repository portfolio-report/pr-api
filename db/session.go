package db

import (
	"time"
)

// Session in database
type Session struct {
	Token          string
	CreatedAt      time.Time
	LastActivityAt time.Time
	Note           string
	UserID         uint
	User           User `gorm:"foreignKey:user_id"`
}

// TableName defines name of table in database
func (Session) TableName() string {
	return "sessions"
}
