package db

type Taxonomy struct {
	UUID        string `gorm:"primaryKey"`
	ParentUUID  *string
	RootUUID    *string
	Name        string
	Code        *string
	Descendants []Taxonomy `gorm:"foreignKey:root_uuid;references:uuid"`
}

func (Taxonomy) TableName() string {
	return "taxonomies"
}
