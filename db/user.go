package db

import (
	"time"

	"gorm.io/datatypes"
)

// User in database
type User struct {
	ID         uint
	Username   string
	Password   *string
	CreatedAt  time.Time
	LastSeenAt datatypes.Date
	IsAdmin    bool
}

// TableName defines name of table in database
func (User) TableName() string {
	return "users"
}
