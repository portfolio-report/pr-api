package models

import "github.com/portfolio-report/pr-api/db"

type ExchangerateResponse struct {
	BaseCurrencyCode  string `json:"baseCurrencyCode"`
	QuoteCurrencyCode string `json:"quoteCurrencyCode"`
}

func ExchangerateResponseFromDB(er db.Exchangerate) ExchangerateResponse {
	return ExchangerateResponse{
		BaseCurrencyCode:  er.BaseCurrencyCode,
		QuoteCurrencyCode: er.QuoteCurrencyCode,
	}
}

type ExchangeratePriceResponse struct {
	Date  db.DbDate `json:"date"`
	Value string    `json:"value"`
}

func ExchangeratePriceResponseFromDB(p db.ExchangeratePrice) ExchangeratePriceResponse {
	return ExchangeratePriceResponse{
		Date:  p.Date,
		Value: *p.Value.String(),
	}
}
