package db

type Security struct {
	UUID               string `gorm:"primaryKey"`
	Name               *string
	Isin               *string
	Wkn                *string
	SymbolXfra         *string
	SymbolXnas         *string
	SymbolXnys         *string
	SecurityType       *string
	SecurityMarkets    []SecurityMarket   `gorm:"foreignKey:security_uuid;references:uuid"`
	Events             []Event            `gorm:"foreignKey:security_uuid;references:uuid"`
	SecurityTaxonomies []SecurityTaxonomy `gorm:"foreignKey:security_uuid;references:uuid"`
}

func (Security) TableName() string {
	return "securities"
}
