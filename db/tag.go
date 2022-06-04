package db

// Tag in database
type Tag struct {
	Name       string     `gorm:"primaryKey"`
	Securities []Security `gorm:"many2many:securities_tags"`
}

// TableName defines name of table in database
func (Tag) TableName() string {
	return "tags"
}
