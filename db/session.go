package db

import (
	"time"
)

type Session struct {
	Token          string
	CreatedAt      time.Time
	LastActivityAt time.Time
	Note           string
	UserID         uint
	User           User `gorm:"foreignKey:user_id"`
}

func (Session) TableName() string {
	return "sessions"
}
