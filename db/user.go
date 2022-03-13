package db

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID         uint
	Username   string
	Password   *string
	CreatedAt  time.Time
	LastSeenAt datatypes.Date
	IsAdmin    bool
}

func (User) TableName() string {
	return "users"
}
