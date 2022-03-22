package db

import "github.com/google/uuid"

// Taxonomy in database
type Taxonomy struct {
	UUID        uuid.UUID `gorm:"primaryKey"`
	ParentUUID  *uuid.UUID
	RootUUID    *uuid.UUID
	Name        string
	Code        *string
	Descendants []Taxonomy `gorm:"foreignKey:root_uuid;references:uuid"`
}

// TableName defines name of table in database
func (Taxonomy) TableName() string {
	return "taxonomies"
}
