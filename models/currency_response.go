package models

import "github.com/portfolio-report/pr-api/db"

type CurrencyResponse struct {
	Code               string                 `json:"code"`
	ExchangeratesBase  []ExchangerateResponse `json:"exchangeratesBase"`
	ExchangeratesQuote []ExchangerateResponse `json:"exchangeratesQuote"`
}

func CurrencyResponseFromDB(c db.Currency) CurrencyResponse {
	exchangeratesBase := []ExchangerateResponse{}
	for _, er := range c.ExchangeratesBase {
		exchangeratesBase = append(exchangeratesBase, ExchangerateResponseFromDB(er))
	}
	exchangeratesQuote := []ExchangerateResponse{}
	for _, er := range c.ExchangeratesQuote {
		exchangeratesQuote = append(exchangeratesQuote, ExchangerateResponseFromDB(er))
	}

	return CurrencyResponse{
		Code:               c.Code,
		ExchangeratesBase:  exchangeratesBase,
		ExchangeratesQuote: exchangeratesQuote,
	}
}
