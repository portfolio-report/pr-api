package db

// Taxonomy in database
type Taxonomy struct {
	UUID        string `gorm:"primaryKey"`
	ParentUUID  *string
	RootUUID    *string
	Name        string
	Code        *string
	Descendants []Taxonomy `gorm:"foreignKey:root_uuid;references:uuid"`
}

// TableName defines name of table in database
func (Taxonomy) TableName() string {
	return "taxonomies"
}
