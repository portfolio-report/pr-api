package db

import "github.com/shopspring/decimal"

// SecurityTaxonomy in database
type SecurityTaxonomy struct {
	SecurityUUID string   `gorm:"primaryKey"`
	TaxonomyUUID string   `gorm:"primaryKey"`
	Taxonomy     Taxonomy `gorm:"foreignKey:taxonomy_uuid"`
	Weight       decimal.Decimal
}

// TableName defines name of table in database
func (SecurityTaxonomy) TableName() string {
	return "securities_taxonomies"
}
