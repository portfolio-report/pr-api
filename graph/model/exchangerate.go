package model

// Exchangerate as used in API
type Exchangerate struct {
	ID                uint   `json:"-"`
	BaseCurrencyCode  string `json:"baseCurrencyCode"`
	QuoteCurrencyCode string `json:"quoteCurrencyCode"`
}
