// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	UUID                 uuid.UUID            `json:"uuid"`
	Type                 PortfolioAccountType `json:"type"`
	Name                 string               `json:"name"`
	CurrencyCode         *string              `json:"currencyCode"`
	ReferenceAccountUUID *uuid.UUID           `json:"referenceAccountUuid"`
	Active               bool                 `json:"active"`
	Note                 string               `json:"note"`
	UpdatedAt            time.Time            `json:"updatedAt"`
	Balance              string               `json:"balance"`
	Value                string               `json:"value"`
}

type PortfolioAccountInput struct {
	Type                 PortfolioAccountType `json:"type"`
	Name                 string               `json:"name"`
	CurrencyCode         *string              `json:"currencyCode"`
	ReferenceAccountUUID *uuid.UUID           `json:"referenceAccountUuid"`
	Active               bool                 `json:"active"`
	Note                 string               `json:"note"`
	UpdatedAt            time.Time            `json:"updatedAt"`
}

type PortfolioInput struct {
	Name             string `json:"name"`
	Note             string `json:"note"`
	BaseCurrencyCode string `json:"baseCurrencyCode"`
}

type PortfolioSecurity struct {
	UUID          uuid.UUID                    `json:"uuid"`
	Name          string                       `json:"name"`
	CurrencyCode  string                       `json:"currencyCode"`
	Isin          string                       `json:"isin"`
	Wkn           string                       `json:"wkn"`
	Symbol        string                       `json:"symbol"`
	Active        bool                         `json:"active"`
	Note          string                       `json:"note"`
	SecurityUUID  *uuid.UUID                   `json:"securityUuid"`
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

type PortfolioSecurityEventInput struct {
	Date    Date   `json:"date"`
	Type    string `json:"type"`
	Details string `json:"details"`
}

type PortfolioSecurityInput struct {
	Name          string                            `json:"name"`
	CurrencyCode  string                            `json:"currencyCode"`
	Isin          string                            `json:"isin"`
	Wkn           string                            `json:"wkn"`
	Symbol        string                            `json:"symbol"`
	Active        bool                              `json:"active"`
	Note          string                            `json:"note"`
	SecurityUUID  *uuid.UUID                        `json:"securityUuid"`
	UpdatedAt     time.Time                         `json:"updatedAt"`
	Calendar      *string                           `json:"calendar"`
	Feed          *string                           `json:"feed"`
	FeedURL       *string                           `json:"feedUrl"`
	LatestFeed    *string                           `json:"latestFeed"`
	LatestFeedURL *string                           `json:"latestFeedUrl"`
	Events        []*PortfolioSecurityEventInput    `json:"events"`
	Properties    []*PortfolioSecurityPropertyInput `json:"properties"`
}

type PortfolioSecurityProperty struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PortfolioSecurityPropertyInput struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PortfolioTransaction struct {
	UUID                   uuid.UUID                   `json:"uuid"`
	AccountUUID            uuid.UUID                   `json:"accountUuid"`
	Type                   PortfolioTransactionType    `json:"type"`
	Datetime               time.Time                   `json:"datetime"`
	PartnerTransactionUUID *uuid.UUID                  `json:"partnerTransactionUuid"`
	Shares                 *decimal.Decimal            `json:"shares"`
	PortfolioSecurityUUID  *uuid.UUID                  `json:"portfolioSecurityUuid"`
	Note                   string                      `json:"note"`
	UpdatedAt              time.Time                   `json:"updatedAt"`
	Units                  []*PortfolioTransactionUnit `json:"units"`
}

type PortfolioTransactionInput struct {
	AccountUUID            uuid.UUID                        `json:"accountUuid"`
	Type                   PortfolioTransactionType         `json:"type"`
	Datetime               time.Time                        `json:"datetime"`
	PartnerTransactionUUID *uuid.UUID                       `json:"partnerTransactionUuid"`
	Shares                 *decimal.Decimal                 `json:"shares"`
	PortfolioSecurityUUID  *uuid.UUID                       `json:"portfolioSecurityUuid"`
	Note                   string                           `json:"note"`
	UpdatedAt              time.Time                        `json:"updatedAt"`
	Units                  []*PortfolioTransactionUnitInput `json:"units"`
}

type PortfolioTransactionUnit struct {
	Type                 PortfolioTransactionUnitType `json:"type"`
	Amount               decimal.Decimal              `json:"amount"`
	CurrencyCode         string                       `json:"currencyCode"`
	OriginalAmount       *decimal.Decimal             `json:"originalAmount"`
	OriginalCurrencyCode *string                      `json:"originalCurrencyCode"`
	ExchangeRate         *decimal.Decimal             `json:"exchangeRate"`
}

type PortfolioTransactionUnitInput struct {
	Type                 PortfolioTransactionUnitType `json:"type"`
	Amount               decimal.Decimal              `json:"amount"`
	CurrencyCode         string                       `json:"currencyCode"`
	OriginalAmount       *decimal.Decimal             `json:"originalAmount"`
	OriginalCurrencyCode *string                      `json:"originalCurrencyCode"`
	ExchangeRate         *decimal.Decimal             `json:"exchangeRate"`
}

type Security struct {
	UUID               uuid.UUID           `json:"uuid"`
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
	SecurityUUID   uuid.UUID `json:"securityUuid"`
	MarketCode     string    `json:"marketCode"`
	CurrencyCode   string    `json:"currencyCode"`
	FirstPriceDate *Date     `json:"firstPriceDate"`
	LastPriceDate  *Date     `json:"lastPriceDate"`
	Symbol         *string   `json:"symbol"`
	UpdatePrices   *bool     `json:"updatePrices"`
}

type SecurityTaxonomy struct {
	SecurityUUID uuid.UUID       `json:"securityUuid"`
	TaxonomyUUID uuid.UUID       `json:"taxonomyUuid"`
	Weight       decimal.Decimal `json:"weight"`
	Taxonomy     *Taxonomy       `json:"taxonomy"`
}

type SecurityTaxonomyInput struct {
	TaxonomyUUID uuid.UUID       `json:"taxonomyUuid"`
	Weight       decimal.Decimal `json:"weight"`
}

type Taxonomy struct {
	UUID       uuid.UUID  `json:"uuid"`
	ParentUUID *uuid.UUID `json:"parentUuid"`
	RootUUID   *uuid.UUID `json:"rootUuid"`
	Name       string     `json:"name"`
	Code       *string    `json:"code"`
}

type TaxonomyInput struct {
	ParentUUID *uuid.UUID `json:"parentUuid"`
	RootUUID   *uuid.UUID `json:"rootUuid"`
	Name       string     `json:"name"`
	Code       *string    `json:"code"`
}

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	IsAdmin    bool   `json:"isAdmin"`
	LastSeenAt string `json:"lastSeenAt"`
}
