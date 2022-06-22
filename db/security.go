package db

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Security in database
type Security struct {
	UUID               uuid.UUID `gorm:"primaryKey"`
	Name               *string
	Isin               *string
	Wkn                *string
	SymbolXfra         *string
	SymbolXnas         *string
	SymbolXnys         *string
	SecurityType       *string
	Extras             datatypes.JSON     `gorm:"default:'{}'"`
	SecurityMarkets    []SecurityMarket   `gorm:"foreignKey:security_uuid;references:uuid"`
	Events             []Event            `gorm:"foreignKey:security_uuid;references:uuid"`
	SecurityTaxonomies []SecurityTaxonomy `gorm:"foreignKey:security_uuid;references:uuid"`
	Tags               []Tag              `gorm:"many2many:securities_tags"`
}

// TableName defines name of table in database
func (Security) TableName() string {
	return "securities"
}
