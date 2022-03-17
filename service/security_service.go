package service

import (
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
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
