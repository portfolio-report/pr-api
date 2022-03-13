package db

import "github.com/shopspring/decimal"

type SecurityTaxonomy struct {
	SecurityUUID string   `gorm:"primaryKey"`
	TaxonomyUUID string   `gorm:"primaryKey"`
	Taxonomy     Taxonomy `gorm:"foreignKey:taxonomy_uuid"`
	Weight       decimal.Decimal
}

func (SecurityTaxonomy) TableName() string {
	return "securities_taxonomies"
}
