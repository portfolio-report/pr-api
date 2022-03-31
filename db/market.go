package db

// Market in database
type Market struct {
	Code string `gorm:"primaryKey"`
	Name string
}

// TableName defines name of table in database
func (Market) TableName() string {
	return "markets"
}
