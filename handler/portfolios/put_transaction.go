package portfolios

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/handler/middleware"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"github.com/shopspring/decimal"
)

type putTransactionUnitRequest struct {
	Type                 string              `json:"type" binding:"oneof=base tax fee"`
	Amount               decimal.Decimal     `json:"amount"`
	CurrencyCode         string              `json:"currencyCode" binding:"max=3"`
	OriginalAmount       decimal.NullDecimal `json:"originalAmount"`
	OriginalCurrencyCode *string             `json:"originalCurrencyCode" binding:"omitempty,max=3"`
	ExchangeRate         decimal.NullDecimal `json:"exchangeRate"`
}

type putTransactionRequest struct {
	AccountUuid            uuid.UUID                   `json:"accountUuid"`
	Type                   string                      `json:"type" binding:"oneof=Payment CurrencyTransfer DepositInterest DepositFee DepositTax SecuritiesOrder SecuritiesDividend SecuritiesFee SecuritiesTax SecuritiesTransfer"`
	Datetime               time.Time                   `json:"datetime"`
	PartnerTransactionUuid *uuid.UUID                  `json:"partnerTransactionUuid"`
	Units                  []putTransactionUnitRequest `json:"units" binding:"dive"`
	Shares                 decimal.NullDecimal         `json:"shares"`
	PortfolioSecurityUuid  *uuid.UUID                  `json:"portfolioSecurityUuid"`
	Note                   string                      `json:"note"`
	UpdatedAt              time.Time                   `json:"updatedAt"`
}

func (h *PortfoliosHandler) PutTransaction(c *gin.Context) {
	portfolioId := uint(middleware.PortfolioFromContext(c).ID)
	uuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		libs.HandleNotFoundError(c)
		return
	}

	var req putTransactionRequest
	if err := c.BindJSON(&req); err != nil {
		libs.HandleBadRequestError(c, err.Error())
		return
	}

	var transaction db.PortfolioTransaction
	err = h.DB.Preload("Units").FirstOrInit(&transaction, db.PortfolioTransaction{PortfolioID: portfolioId, UUID: uuid}).Error
	if err != nil {
		panic(err)
	}

	transaction.Type = req.Type
	transaction.Datetime = req.Datetime
	transaction.Note = req.Note
	transaction.Shares = req.Shares
	transaction.UpdatedAt = req.UpdatedAt
	transaction.AccountUUID = req.AccountUuid
	transaction.PartnerTransactionUUID = req.PartnerTransactionUuid
	transaction.PortfolioSecurityUUID = req.PortfolioSecurityUuid

	if err := h.DB.Save(&transaction).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			libs.HandleBadRequestError(c, "data violates constraint "+pqErr.Constraint)
			return
		}

		panic(err)
	}

	units, err := h.createUpdateDeleteTransactionUnits(portfolioId, uuid, req.Units, transaction.Units)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			libs.HandleBadRequestError(c, "data violates constraint "+pqErr.Constraint)
			return
		}

		panic(err)
	}

	transaction.Units = units

	c.JSON(http.StatusOK, models.PortfolioTransactionResponseFromDB(&transaction))
}

// Create/update/delete units in database to match units in request
func (h *PortfoliosHandler) createUpdateDeleteTransactionUnits(
	portfolioId uint,
	transactionUuid uuid.UUID,
	req []putTransactionUnitRequest,
	dbUnits []db.PortfolioTransactionUnit,
) ([]db.PortfolioTransactionUnit, error) {

	matcher := func(r putTransactionUnitRequest, dbUnit db.PortfolioTransactionUnit) bool {
		return dbUnit.Type == r.Type &&
			dbUnit.Amount.Equal(r.Amount) &&
			dbUnit.CurrencyCode == r.CurrencyCode &&
			equalNullDecimal(dbUnit.OriginalAmount, r.OriginalAmount) &&
			((dbUnit.OriginalCurrencyCode == nil && r.OriginalCurrencyCode == nil) || *dbUnit.OriginalCurrencyCode == *r.OriginalCurrencyCode) &&
			equalNullDecimal(dbUnit.ExchangeRate, r.ExchangeRate)
	}

	unmatchedReq, unmatchedDb, matchedDb := libs.MatchElementsInArrays(req, dbUnits, matcher)

	// Delete removed units
	if len(unmatchedDb) > 0 {
		err := h.DB.Delete(unmatchedDb).Error
		if err != nil {
			panic(err)
		}
	}

	// Create new units
	newDb := []db.PortfolioTransactionUnit{}
	for _, el := range unmatchedReq {
		newDb = append(newDb, db.PortfolioTransactionUnit{
			PortfolioID:          portfolioId,
			TransactionUUID:      transactionUuid,
			Type:                 el.Type,
			Amount:               el.Amount,
			CurrencyCode:         el.CurrencyCode,
			OriginalAmount:       el.OriginalAmount,
			OriginalCurrencyCode: el.OriginalCurrencyCode,
			ExchangeRate:         el.ExchangeRate,
		})
	}

	if len(newDb) > 0 {
		if err := h.DB.Create(&newDb).Error; err != nil {
			return nil, err
		}
	}

	return append(matchedDb, newDb...), nil
}

func equalNullDecimal(d1 decimal.NullDecimal, d2 decimal.NullDecimal) bool {
	if !d1.Valid && !d2.Valid {
		return true
	}
	if d1.Valid && d2.Valid {
		return d1.Decimal.Equal(d2.Decimal)
	}
	return false
}
