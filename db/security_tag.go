package db

import "github.com/google/uuid"

// SecurityTag in database
type SecurityTag struct {
	SecurityUUID uuid.UUID `gorm:"primaryKey"`
	TagName      string    `gorm:"primaryKey"`
}

// TableName defines name of table in database
func (SecurityTag) TableName() string {
	return "securities_tags"
}
