package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/db"
	"github.com/shopspring/decimal"
)

type PortfolioTransactionResponse struct {
	UUID                   uuid.UUID
	AccountUUID            uuid.UUID
	Type                   string
	Datetime               time.Time
	PartnerTransactionUUID *uuid.UUID
	Shares                 decimal.NullDecimal
	PortfolioSecurityUUID  *uuid.UUID
	Note                   string
	UpdatedAt              time.Time
	Units                  []PortfolioTransactionUnitResponse
}

func PortfolioTransactionResponseFromDB(t *db.PortfolioTransaction) PortfolioTransactionResponse {
	units := []PortfolioTransactionUnitResponse{}
	for _, u := range t.Units {
		units = append(units, PortfolioTransactionUnitResponseFromDB(&u))
	}
	return PortfolioTransactionResponse{
		UUID:                   t.UUID,
		AccountUUID:            t.AccountUUID,
		Type:                   t.Type,
		Datetime:               t.Datetime.UTC(),
		PartnerTransactionUUID: t.PartnerTransactionUUID,
		Shares:                 t.Shares,
		PortfolioSecurityUUID:  t.PortfolioSecurityUUID,
		Note:                   t.Note,
		UpdatedAt:              t.UpdatedAt.UTC(),
		Units:                  units,
	}
}

type PortfolioTransactionUnitResponse struct {
	Type                 string              `json:"type"`
	Amount               string              `json:"amount"`
	CurrencyCode         string              `json:"currencyCode"`
	OriginalAmount       decimal.NullDecimal `json:"originalAmount"`
	OriginalCurrencyCode *string             `json:"originalCurrencyCode"`
	ExchangeRate         decimal.NullDecimal `json:"exchangeRate"`
}

func PortfolioTransactionUnitResponseFromDB(u *db.PortfolioTransactionUnit) PortfolioTransactionUnitResponse {
	return PortfolioTransactionUnitResponse{
		Type:                 u.Type,
		Amount:               u.Amount.String(),
		CurrencyCode:         u.CurrencyCode,
		OriginalAmount:       u.OriginalAmount,
		OriginalCurrencyCode: u.OriginalCurrencyCode,
		ExchangeRate:         u.ExchangeRate,
	}
}
