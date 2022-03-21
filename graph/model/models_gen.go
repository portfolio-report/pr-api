// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Currency struct {
	Code               string          `json:"code"`
	ExchangeratesBase  []*Exchangerate `json:"exchangeratesBase"`
	ExchangeratesQuote []*Exchangerate `json:"exchangeratesQuote"`
}

type Event struct {
	Date         string  `json:"date"`
	Type         string  `json:"type"`
	Amount       *string `json:"amount"`
	CurrencyCode *string `json:"currencyCode"`
	Ratio        *string `json:"ratio"`
}

type ExchangeratePrice struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

type Portfolio struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Note             string    `json:"note"`
	BaseCurrencyCode string    `json:"baseCurrencyCode"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type PortfolioAccount struct {
	UUID                 string    `json:"uuid"`
	Type                 string    `json:"type"`
	Name                 string    `json:"name"`
	CurrencyCode         *string   `json:"currencyCode"`
	ReferenceAccountUUID *string   `json:"referenceAccountUuid"`
	Active               bool      `json:"active"`
	Note                 string    `json:"note"`
	UpdatedAt            time.Time `json:"updatedAt"`
	Balance              string    `json:"balance"`
	Value                string    `json:"value"`
}

type PortfolioInput struct {
	Name             string `json:"name"`
	Note             string `json:"note"`
	BaseCurrencyCode string `json:"baseCurrencyCode"`
}

type PortfolioSecurity struct {
	UUID          string                       `json:"uuid"`
	Name          string                       `json:"name"`
	CurrencyCode  string                       `json:"currencyCode"`
	Isin          string                       `json:"isin"`
	Wkn           string                       `json:"wkn"`
	Symbol        string                       `json:"symbol"`
	Active        bool                         `json:"active"`
	Note          string                       `json:"note"`
	SecurityUUID  *string                      `json:"securityUuid"`
	UpdatedAt     time.Time                    `json:"updatedAt"`
	Calendar      *string                      `json:"calendar"`
	Feed          *string                      `json:"feed"`
	FeedURL       *string                      `json:"feedUrl"`
	LatestFeed    *string                      `json:"latestFeed"`
	LatestFeedURL *string                      `json:"latestFeedUrl"`
	Events        []*PortfolioSecurityEvent    `json:"events"`
	Properties    []*PortfolioSecurityProperty `json:"properties"`
	Shares        string                       `json:"shares"`
	Quote         string                       `json:"quote"`
}

type PortfolioSecurityEvent struct {
	Date    Date   `json:"date"`
	Type    string `json:"type"`
	Details string `json:"details"`
}

type PortfolioSecurityProperty struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Security struct {
	UUID               string              `json:"uuid"`
	Name               *string             `json:"name"`
	Isin               *string             `json:"isin"`
	Wkn                *string             `json:"wkn"`
	SecurityType       *string             `json:"securityType"`
	SymbolXfra         *string             `json:"symbolXfra"`
	SymbolXnas         *string             `json:"symbolXnas"`
	SymbolXnys         *string             `json:"symbolXnys"`
	SecurityMarkets    []*SecurityMarket   `json:"securityMarkets"`
	SecurityTaxonomies []*SecurityTaxonomy `json:"securityTaxonomies"`
	Events             []*Event            `json:"events"`
}

type SecurityInput struct {
	Name         *string `json:"name"`
	Isin         *string `json:"isin"`
	Wkn          *string `json:"wkn"`
	SecurityType *string `json:"securityType"`
	SymbolXfra   *string `json:"symbolXfra"`
	SymbolXnas   *string `json:"symbolXnas"`
	SymbolXnys   *string `json:"symbolXnys"`
}

type SecurityMarket struct {
	SecurityUUID   string  `json:"securityUuid"`
	MarketCode     string  `json:"marketCode"`
	CurrencyCode   string  `json:"currencyCode"`
	FirstPriceDate *string `json:"firstPriceDate"`
	LastPriceDate  *string `json:"lastPriceDate"`
	Symbol         *string `json:"symbol"`
	UpdatePrices   *bool   `json:"updatePrices"`
}

type SecurityTaxonomy struct {
	SecurityUUID string    `json:"securityUuid"`
	TaxonomyUUID string    `json:"taxonomyUuid"`
	Weight       string    `json:"weight"`
	Taxonomy     *Taxonomy `json:"taxonomy"`
}

type Taxonomy struct {
	UUID       string  `json:"uuid"`
	ParentUUID *string `json:"parentUuid"`
	RootUUID   *string `json:"rootUuid"`
	Name       string  `json:"name"`
	Code       *string `json:"code"`
}

type TaxonomyInput struct {
	ParentUUID *string `json:"parentUuid"`
	RootUUID   *string `json:"rootUuid"`
	Name       *string `json:"name"`
	Code       *string `json:"code"`
}

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	IsAdmin    bool   `json:"isAdmin"`
	LastSeenAt string `json:"lastSeenAt"`
}
