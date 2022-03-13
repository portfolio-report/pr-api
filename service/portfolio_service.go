package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type portfolioService struct {
	DB *gorm.DB
}

func NewPortfolioService(db *gorm.DB) models.PortfolioService {
	return &portfolioService{
		DB: db,
	}
}

func (*portfolioService) ModelFromDb(p db.Portfolio) *model.Portfolio {
	return &model.Portfolio{
		ID:               int(p.ID),
		Name:             p.Name,
		BaseCurrencyCode: p.BaseCurrencyCode,
		Note:             p.Note,
		CreatedAt:        p.CreatedAt.UTC(),
		UpdatedAt:        p.UpdatedAt.UTC(),
	}
}

func (s *portfolioService) GetAllOfUser(user *model.User) []*model.Portfolio {
	var portfolios []db.Portfolio
	err := s.DB.Find(&portfolios, "user_id = ?", user.ID).Error
	if err != nil {
		panic(err)
	}

	response := []*model.Portfolio{}
	for _, p := range portfolios {
		response = append(response, s.ModelFromDb(p))
	}
	return response
}

func (s *portfolioService) GetPortfolioOfUserByID(user *model.User, ID uint) (*model.Portfolio, error) {
	var portfolio db.Portfolio
	if err := s.DB.Take(&portfolio, "user_id = ? AND id = ?", user.ID, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		panic(err)
	}
	return s.ModelFromDb(portfolio), nil
}

func (s *portfolioService) GetPortfolioByID(ID uint) (*model.Portfolio, error) {

	var portfolio db.Portfolio
	if err := s.DB.Take(&portfolio, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		panic(err)
	}
	return s.ModelFromDb(portfolio), nil
}

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

	return s.ModelFromDb(portfolio), nil
}

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

	return s.ModelFromDb(portfolio), nil
}

func (s *portfolioService) DeletePortfolio(ID uint) *model.Portfolio {
	portfolio := db.Portfolio{
		ID: ID,
	}
	err := s.DB.Clauses(clause.Returning{}).Delete(&portfolio).Error
	if err != nil {
		panic(err)
	}
	return s.ModelFromDb(portfolio)
}
