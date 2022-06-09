package db

import "github.com/google/uuid"

// SecurityTag in database
type SecurityTag struct {
	SecurityUUID uuid.UUID `gorm:"primaryKey"`
	TagUUID      uuid.UUID `gorm:"primaryKey"`
	Tag          Tag       `gorm:"foreignKey:tag_uuid"`
}

// TableName defines name of table in database
func (SecurityTag) TableName() string {
	return "securities_tags"
}
