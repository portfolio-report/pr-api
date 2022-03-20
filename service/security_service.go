package service

import (
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type securityService struct {
	DB *gorm.DB
}

// NewSecurityService creates and returns new security service
func NewSecurityService(db *gorm.DB) models.SecurityService {
	return &securityService{
		DB: db,
	}
}

// GetSecurityByUUID returns security idenfitied by UUID
func (s *securityService) GetSecurityByUUID(uuid string) (*model.Security, error) {
	var security db.Security
	if err := s.DB.Take(&security, "uuid = ?", uuid).Error; err != nil {
		return nil, err
	}
	return s.modelFromDb(security), nil
}

// GetEventsOfSecurity returns all events of security
func (s *securityService) GetEventsOfSecurity(security *model.Security) ([]*model.Event, error) {
	var events []db.Event
	err := s.DB.Find(&events, "security_uuid = ?", security.UUID).Error
	if err != nil {
		panic(err)
	}
	return s.eventsModelFromDb(events), nil
}

// CreateSecurity create security
func (s *securityService) CreateSecurity(input *model.SecurityInput) (*model.Security, error) {
	security := db.Security{
		UUID:         uuid.New().String(),
		Name:         input.Name,
		Isin:         input.Isin,
		Wkn:          input.Wkn,
		SecurityType: input.SecurityType,
		SymbolXfra:   input.SymbolXfra,
		SymbolXnas:   input.SymbolXnas,
		SymbolXnys:   input.SymbolXnys,
	}

	err := s.DB.Clauses(clause.Returning{}).Create(&security).Error
	if err != nil {
		panic(err)
	}

	return s.modelFromDb(security), nil
}

// UpdateSecurity stores all attributes of input (incl. nil values)
func (s *securityService) UpdateSecurity(uuid string, input *model.SecurityInput) (*model.Security, error) {
	security := db.Security{UUID: uuid}
	err := s.DB.Model(&security).
		Clauses(clause.Returning{}).
		Updates(map[string]interface{}{
			"Name":         input.Name,
			"Isin":         input.Isin,
			"Wkn":          input.Wkn,
			"SecurityType": input.SecurityType,
			"SymbolXfra":   input.SymbolXfra,
			"SymbolXnas":   input.SymbolXnas,
			"SymbolXnys":   input.SymbolXnys,
		}).Error
	if err != nil {
		panic(err)
	}

	return s.modelFromDb(security), nil
}

// DeleteSecurity removes security
func (s *securityService) DeleteSecurity(uuid string) (*model.Security, error) {
	var security db.Security
	result := s.DB.Clauses(clause.Returning{}).Delete(&security, "uuid = ?", uuid)
	err := result.Error
	if err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		return nil, model.ErrNotFound
	}
	return s.modelFromDb(security), nil
}

// eventsModelFromDb converts list of events from database into model
func (*securityService) eventsModelFromDb(events []db.Event) []*model.Event {
	ret := []*model.Event{}
	for _, e := range events {
		ret = append(ret, &model.Event{
			Date:         e.Date.String(),
			Type:         e.Type,
			Amount:       e.Amount,
			CurrencyCode: e.CurrencyCode,
			Ratio:        e.Ratio,
		})
	}
	return ret
}

// modelFromDb converts security from database into model
func (*securityService) modelFromDb(s db.Security) *model.Security {
	return &model.Security{
		UUID:         s.UUID,
		Name:         s.Name,
		Isin:         s.Isin,
		Wkn:          s.Wkn,
		SecurityType: s.SecurityType,
		SymbolXfra:   s.SymbolXfra,
		SymbolXnas:   s.SymbolXnas,
		SymbolXnys:   s.SymbolXnys,
	}
}
