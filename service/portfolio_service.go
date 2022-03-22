package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type portfolioService struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

// NewPortfolioService creates and returns new portfolio service
func NewPortfolioService(db *gorm.DB, validate *validator.Validate) models.PortfolioService {
	return &portfolioService{
		DB:       db,
		Validate: validate,
	}
}

// modelFromDb converts portfolio from database into model
func (*portfolioService) modelFromDb(p db.Portfolio) *model.Portfolio {
	return &model.Portfolio{
		ID:               int(p.ID),
		Name:             p.Name,
		BaseCurrencyCode: p.BaseCurrencyCode,
		Note:             p.Note,
		CreatedAt:        p.CreatedAt.UTC(),
		UpdatedAt:        p.UpdatedAt.UTC(),
	}
}

// accountModelFromDb converts portfolio account from database into model
func (*portfolioService) accountModelFromDb(a db.PortfolioAccount) *model.PortfolioAccount {
	return &model.PortfolioAccount{
		UUID:                 a.UUID,
		Type:                 a.Type,
		Name:                 a.Name,
		CurrencyCode:         a.CurrencyCode,
		ReferenceAccountUUID: a.ReferenceAccountUUID,
		Active:               a.Active,
		Note:                 a.Note,
		UpdatedAt:            a.UpdatedAt.UTC(),
	}
}

// securityModelFromDb converts portfolio security from database into model
func (*portfolioService) securityModelFromDb(s db.PortfolioSecurity) *model.PortfolioSecurity {
	var properties []model.PortfolioSecurityProperty
	err := json.Unmarshal(s.Properties, &properties)
	if err != nil {
		panic(err)
	}
	propertiesPtr := make([]*model.PortfolioSecurityProperty, len(properties))
	for i := range properties {
		propertiesPtr[i] = &properties[i]
	}

	var events []model.PortfolioSecurityEvent
	err = json.Unmarshal(s.Events, &events)
	if err != nil {
		panic(err)
	}
	eventsPtr := make([]*model.PortfolioSecurityEvent, len(events))
	for i := range events {
		eventsPtr[i] = &events[i]
	}

	return &model.PortfolioSecurity{
		UUID:          s.UUID,
		Name:          s.Name,
		CurrencyCode:  s.CurrencyCode,
		Isin:          s.Isin,
		Wkn:           s.Wkn,
		Symbol:        s.Symbol,
		Active:        s.Active,
		Note:          s.Note,
		SecurityUUID:  s.SecurityUUID,
		UpdatedAt:     s.UpdatedAt.UTC(),
		Calendar:      s.Calendar,
		Feed:          s.Feed,
		FeedURL:       s.FeedUrl,
		LatestFeed:    s.LatestFeed,
		LatestFeedURL: s.LatestFeedUrl,
		Events:        eventsPtr,
		Properties:    propertiesPtr,
	}
}

// transactionModelFromDb converts portfolio transaction from database into model
func (*portfolioService) transactionModelFromDb(t db.PortfolioTransaction) *model.PortfolioTransaction {
	units := []*model.PortfolioTransactionUnit{}
	for _, u := range t.Units {
		units = append(units, &model.PortfolioTransactionUnit{
			Type:                 u.Type,
			Amount:               u.Amount,
			CurrencyCode:         u.CurrencyCode,
			OriginalAmount:       u.OriginalAmount,
			OriginalCurrencyCode: u.OriginalCurrencyCode,
			ExchangeRate:         u.ExchangeRate,
		})
	}
	return &model.PortfolioTransaction{
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

// GetAllOfUser returns all portfolios of user
func (s *portfolioService) GetAllOfUser(user *model.User) ([]*model.Portfolio, error) {
	var portfolios []db.Portfolio
	err := s.DB.Find(&portfolios, "user_id = ?", user.ID).Error
	if err != nil {
		panic(err)
	}

	response := []*model.Portfolio{}
	for _, p := range portfolios {
		response = append(response, s.modelFromDb(p))
	}
	return response, nil
}

// GetPortfolioOfUserByID returns single portfolio of user
func (s *portfolioService) GetPortfolioOfUserByID(user *model.User, ID uint) (*model.Portfolio, error) {
	var portfolio db.Portfolio
	if err := s.DB.Take(&portfolio, "user_id = ? AND id = ?", user.ID, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		panic(err)
	}
	return s.modelFromDb(portfolio), nil
}

// GetPortfolioByID returns single portfolio
func (s *portfolioService) GetPortfolioByID(ID uint) (*model.Portfolio, error) {

	var portfolio db.Portfolio
	if err := s.DB.Take(&portfolio, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		panic(err)
	}
	return s.modelFromDb(portfolio), nil
}

// CreatePortfolio creates new portfolio
func (s *portfolioService) CreatePortfolio(user *model.User, req *model.PortfolioInput) (*model.Portfolio, error) {
	portfolio := db.Portfolio{
		Name:             req.Name,
		Note:             req.Note,
		BaseCurrencyCode: req.BaseCurrencyCode,
		UserID:           uint(user.ID),
	}

	err := s.DB.Clauses(clause.Returning{}).Create(&portfolio).Error
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates constraint " + pqErr.Constraint)
		}

		panic(err)
	}

	return s.modelFromDb(portfolio), nil
}

// UpdatePortfolio updates portfolio
func (s *portfolioService) UpdatePortfolio(ID uint, req *model.PortfolioInput) (*model.Portfolio, error) {
	portfolio := db.Portfolio{
		ID:               ID,
		Name:             req.Name,
		Note:             req.Note,
		BaseCurrencyCode: req.BaseCurrencyCode,
		UpdatedAt:        time.Now(),
	}

	err := s.DB.Clauses(clause.Returning{}).Updates(&portfolio).Error
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates contraint " + pqErr.Constraint)
		}

		panic(err)
	}

	return s.modelFromDb(portfolio), nil
}

// DeletePortfolio removes portfolio
func (s *portfolioService) DeletePortfolio(ID uint) (*model.Portfolio, error) {
	portfolio := db.Portfolio{
		ID: ID,
	}
	err := s.DB.Clauses(clause.Returning{}).Delete(&portfolio).Error
	if err != nil {
		panic(err)
	}
	return s.modelFromDb(portfolio), nil
}

// GetPortfolioAccountsOfPortfolio lists all account in portfolio
func (s *portfolioService) GetPortfolioAccountsOfPortfolio(portfolioId int) ([]*model.PortfolioAccount, error) {
	var accounts []db.PortfolioAccount
	err := s.DB.Where("portfolio_id = ?", portfolioId).Find(&accounts).Error
	if err != nil {
		panic(err)
	}

	response := make([]*model.PortfolioAccount, len(accounts))
	for i := range accounts {
		response[i] = s.accountModelFromDb(accounts[i])
	}

	return response, nil
}

// UpsertPortfolioAccount creates or updates portfolio account
func (s *portfolioService) UpsertPortfolioAccount(
	portfolioId int, uuid uuid.UUID, input model.PortfolioAccountInput,
) (
	*model.PortfolioAccount, error,
) {
	var account db.PortfolioAccount

	err := s.DB.FirstOrInit(&account, db.PortfolioAccount{PortfolioID: uint(portfolioId), UUID: uuid}).Error
	if err != nil {
		panic(err)
	}

	account.Type = input.Type
	account.Name = input.Name
	account.Active = input.Active
	account.Note = input.Note
	account.UpdatedAt = input.UpdatedAt

	switch input.Type {
	case "deposit":
		if input.CurrencyCode == nil {
			return nil, fmt.Errorf("currencyCode is missing")
		}
		account.CurrencyCode = input.CurrencyCode
		account.ReferenceAccountUUID = nil
	case "securities":
		if input.ReferenceAccountUUID == nil {
			return nil, fmt.Errorf("referenceAccountUuid is missing")
		}
		account.CurrencyCode = nil
		account.ReferenceAccountUUID = input.ReferenceAccountUUID
	default:
		return nil, fmt.Errorf("invalid type: %s", input.Type)
	}

	if err := s.DB.Save(&account).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates constraint: %s", pqErr.Constraint)
		}

		panic(err)
	}

	return s.accountModelFromDb(account), nil
}

// DeletePortfolioAccount removes account from portfolio and links to it
func (s *portfolioService) DeletePortfolioAccount(portfolioId int, uuid uuid.UUID) (*model.PortfolioAccount, error) {
	// Remove links as reference account
	err := s.DB.Model(&db.PortfolioAccount{}).
		Where("portfolio_id = ? AND reference_account_uuid = ?", portfolioId, uuid).
		Update("reference_account_uuid", nil).Error
	if err != nil {
		panic(err)
	}

	// Delete transactions of account
	err = s.DB.
		Where("portfolio_id = ? AND account_uuid = ?", portfolioId, uuid).
		Delete(&db.PortfolioTransaction{}).Error
	if err != nil {
		panic(err)
	}

	var account db.PortfolioAccount
	result := s.DB.Clauses(clause.Returning{}).Where("portfolio_id = ? AND uuid = ?", portfolioId, uuid).Delete(&account)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		return nil, model.ErrNotFound
	}

	return s.accountModelFromDb(account), nil
}

// GetPortfolioSecuritiesOfPortfolio lists all securities in portfolio
func (s *portfolioService) GetPortfolioSecuritiesOfPortfolio(portfolioId int) ([]*model.PortfolioSecurity, error) {
	var securities []db.PortfolioSecurity
	err := s.DB.Where("portfolio_id = ?", portfolioId).Find(&securities).Error
	if err != nil {
		panic(err)
	}

	response := make([]*model.PortfolioSecurity, len(securities))
	for i := range securities {
		response[i] = s.securityModelFromDb(securities[i])
	}

	return response, nil
}

// UpsertPortfolioSecurity creates or updates portfolio security
func (s *portfolioService) UpsertPortfolioSecurity(
	portfolioId int, uuid uuid.UUID, input model.PortfolioSecurityInput,
) (
	*model.PortfolioSecurity, error,
) {
	var security db.PortfolioSecurity

	err := s.DB.FirstOrInit(&security, db.PortfolioSecurity{PortfolioID: uint(portfolioId), UUID: uuid}).Error
	if err != nil {
		panic(err)
	}

	security.Name = input.Name
	security.CurrencyCode = input.CurrencyCode
	security.Isin = input.Isin
	security.Wkn = input.Wkn
	security.Symbol = input.Symbol
	security.Active = input.Active
	security.Note = input.Note
	security.SecurityUUID = input.SecurityUUID
	security.UpdatedAt = input.UpdatedAt
	security.Calendar = input.Calendar
	security.Feed = input.Feed
	security.FeedUrl = input.FeedURL
	security.LatestFeed = input.LatestFeed
	security.LatestFeedUrl = input.LatestFeedURL
	for _, e := range input.Events {
		if e.Type != "STOCK_SPLIT" && e.Type != "NOTE" && e.Type != "DIVIDEND_PAYMENT" {
			return nil, fmt.Errorf("event type %s is not supported", e.Type)
		}
	}
	security.Events, err = json.Marshal(input.Events)
	if err != nil {
		panic(err)
	}
	for _, p := range input.Properties {
		if p.Type != "MARKET" && p.Type != "FEED" {
			return nil, fmt.Errorf("propertey type %s is not supported", p.Type)
		}
	}
	security.Properties, err = json.Marshal(input.Properties)
	if err != nil {
		panic(err)
	}

	if err := s.DB.Save(&security).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates constraint: %s", pqErr.Constraint)
		}

		panic(err)
	}

	return s.securityModelFromDb(security), nil
}

// DeletePortfolioSecurity removes security from portfolio and links to it
func (s *portfolioService) DeletePortfolioSecurity(portfolioId int, uuid uuid.UUID) (*model.PortfolioSecurity, error) {
	// Delete transactions of security
	err := s.DB.
		Where("portfolio_id = ? AND portfolio_security_uuid = ?", portfolioId, uuid).
		Delete(&db.PortfolioTransaction{}).Error
	if err != nil {
		panic(err)
	}

	var security db.PortfolioSecurity
	result := s.DB.
		Clauses(clause.Returning{}).
		Where("portfolio_id = ? AND uuid = ?", portfolioId, uuid).
		Delete(&security)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		return nil, model.ErrNotFound
	}

	return s.securityModelFromDb(security), nil
}

// GetPortfolioTransactionsOfPortfolio lists all transactions in portfolio
func (s *portfolioService) GetPortfolioTransactionsOfPortfolio(portfolioId int) ([]*model.PortfolioTransaction, error) {
	var transactions []db.PortfolioTransaction
	err := s.DB.
		Preload("Units").
		Where("portfolio_id = ?", portfolioId).
		Find(&transactions).Error
	if err != nil {
		panic(err)
	}

	response := make([]*model.PortfolioTransaction, len(transactions))
	for i := range transactions {
		response[i] = s.transactionModelFromDb(transactions[i])
	}

	return response, nil
}

// UpsertPortfolioTransaction creates or updates portfolio transaction
func (s *portfolioService) UpsertPortfolioTransaction(
	portfolioId int, uuid uuid.UUID, input model.PortfolioTransactionInput,
) (
	*model.PortfolioTransaction, error,
) {
	var transaction db.PortfolioTransaction
	err := s.DB.
		Preload("Units").
		FirstOrInit(&transaction, db.PortfolioTransaction{PortfolioID: uint(portfolioId), UUID: uuid}).Error
	if err != nil {
		panic(err)
	}

	transaction.Type = input.Type
	transaction.Datetime = input.Datetime
	transaction.Note = input.Note
	transaction.Shares = input.Shares
	transaction.UpdatedAt = input.UpdatedAt
	transaction.AccountUUID = input.AccountUUID
	transaction.PartnerTransactionUUID = input.PartnerTransactionUUID
	transaction.PortfolioSecurityUUID = input.PortfolioSecurityUUID

	if err := s.DB.Save(&transaction).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates constraint %s", pqErr.Constraint)
		}

		panic(err)
	}

	units, err := s.createUpdateDeleteTransactionUnits(uint(portfolioId), uuid, input.Units, transaction.Units)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates constraint %s", pqErr.Constraint)
		}

		panic(err)
	}

	transaction.Units = units

	return s.transactionModelFromDb(transaction), nil
}

// createUpdateDeleteTransactionUnits creates/updates/deletes units in database to match units in request
func (s *portfolioService) createUpdateDeleteTransactionUnits(
	portfolioId uint,
	transactionUuid uuid.UUID,
	req []*model.PortfolioTransactionUnitInput,
	dbUnits []db.PortfolioTransactionUnit,
) ([]db.PortfolioTransactionUnit, error) {

	matcher := func(r *model.PortfolioTransactionUnitInput, dbUnit db.PortfolioTransactionUnit) bool {
		return dbUnit.Type == r.Type &&
			dbUnit.Amount.Equal(r.Amount) &&
			dbUnit.CurrencyCode == r.CurrencyCode &&
			equalDecimalPtr(dbUnit.OriginalAmount, r.OriginalAmount) &&
			((dbUnit.OriginalCurrencyCode == nil && r.OriginalCurrencyCode == nil) || *dbUnit.OriginalCurrencyCode == *r.OriginalCurrencyCode) &&
			equalDecimalPtr(dbUnit.ExchangeRate, r.ExchangeRate)
	}

	unmatchedReq, unmatchedDb, matchedDb := libs.MatchElementsInArrays(req, dbUnits, matcher)

	// Delete removed units
	if len(unmatchedDb) > 0 {
		err := s.DB.Delete(unmatchedDb).Error
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
		if err := s.DB.Create(&newDb).Error; err != nil {
			return nil, err
		}
	}

	return append(matchedDb, newDb...), nil
}

func equalDecimalPtr(d1 *decimal.Decimal, d2 *decimal.Decimal) bool {
	if d1 == nil && d2 == nil {
		return true
	}
	if d1 != nil && d2 != nil {
		return d1.Equal(*d2)
	}
	return false
}

// DeletePortfolioTransaction removes transaction from portfolio and links to it
func (s *portfolioService) DeletePortfolioTransaction(portfolioId int, uuid uuid.UUID) (*model.PortfolioTransaction, error) {
	// Remove link from partner transaction (if exists)
	s.DB.Model(&db.PortfolioTransaction{}).
		Where("portfolio_id = ? AND partner_transaction_uuid = ?", portfolioId, uuid).
		Update("partner_transaction_uuid", nil)

	var transaction db.PortfolioTransaction
	result := s.DB.
		Clauses(clause.Returning{}).
		Where("portfolio_id = ? AND uuid = ?", portfolioId, uuid).
		Delete(&transaction)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		return nil, model.ErrNotFound
	}

	return s.transactionModelFromDb(transaction), nil
}
