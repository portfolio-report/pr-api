package models

import "github.com/portfolio-report/pr-api/db"

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
