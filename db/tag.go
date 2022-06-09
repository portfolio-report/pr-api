package db

import "github.com/google/uuid"

// Tag in database
type Tag struct {
	UUID       uuid.UUID `gorm:"primaryKey"`
	Name       string
	Securities []Security `gorm:"many2many:securities_tags"`
}

// TableName defines name of table in database
func (Tag) TableName() string {
	return "tags"
}
