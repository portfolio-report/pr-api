package models

import (
	"time"

	"github.com/portfolio-report/pr-api/db"
)

type PortfolioAccountResponse struct {
	UUID                 string    `json:"uuid"`
	Type                 string    `json:"type"`
	Name                 string    `json:"name"`
	CurrencyCode         *string   `json:"currencyCode"`
	ReferenceAccountUUID *string   `json:"referenceAccountUuid"`
	Active               bool      `json:"active"`
	Note                 string    `json:"note"`
	UpdatedAt            time.Time `json:"updateAt"`
}

func PortfolioAccountResponseFromDB(e *db.PortfolioAccount) PortfolioAccountResponse {
	return PortfolioAccountResponse{
		UUID:                 e.UUID,
		Type:                 e.Type,
		Name:                 e.Name,
		CurrencyCode:         e.CurrencyCode,
		ReferenceAccountUUID: e.ReferenceAccountUUID,
		Active:               e.Active,
		Note:                 e.Note,
		UpdatedAt:            e.UpdatedAt.UTC(),
	}
}
