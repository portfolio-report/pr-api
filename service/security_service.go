package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type securityService struct {
	DB *gorm.DB
}

// NewSecurityService creates and returns new security service
func NewSecurityService(db *gorm.DB) model.SecurityService {
	return &securityService{
		DB: db,
	}
}

// GetSecurityByUUID returns security idenfitied by UUID
func (s *securityService) GetSecurityByUUID(uuid uuid.UUID) (*model.Security, error) {
	var security db.Security
	if err := s.DB.Take(&security, "uuid = ?", uuid).Error; err != nil {
		return nil, err
	}
	return s.modelFromDb(security), nil
}

// GetSecuritiesByTag returns securities associated with tag
func (s *securityService) GetSecuritiesByTag(tag string) []*model.Security {
	var dbTag db.Tag
	if err := s.DB.Take(&dbTag, "LOWER(name) = LOWER(?)", tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*model.Security{}
		}
		panic(err)
	}

	var dbSecurities []db.Security
	if err := s.DB.Model(&dbTag).Association("Securities").Find(&dbSecurities); err != nil {
		panic(err)
	}

	securities := make([]*model.Security, len(dbSecurities))
	for i := range dbSecurities {
		securities[i] = s.modelFromDb(dbSecurities[i])
	}

	return securities
}

// GetEventsOfSecurity returns all events of security
func (s *securityService) GetEventsOfSecurity(security *model.Security) []*model.Event {
	var events []db.Event
	err := s.DB.Find(&events, "security_uuid = ?", security.UUID).Error
	if err != nil {
		panic(err)
	}
	return s.eventsModelFromDb(events)
}

// CreateSecurity create security
func (s *securityService) CreateSecurity(input *model.SecurityInput) (*model.Security, error) {
	security := db.Security{
		UUID:         uuid.New(),
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
func (s *securityService) UpdateSecurity(uuid uuid.UUID, input *model.SecurityInput) (*model.Security, error) {
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
func (s *securityService) DeleteSecurity(uuid uuid.UUID) (*model.Security, error) {
	var security db.Security
	result := s.DB.Clauses(clause.Returning{}).Delete(&security, "uuid = ?", uuid)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		return nil, model.ErrNotFound
	}
	return s.modelFromDb(security), nil
}

// DeleteSecurityMarket removes market of security
func (s *securityService) DeleteSecurityMarket(securityUuid uuid.UUID, marketCode string) (*model.SecurityMarket, error) {
	var market db.SecurityMarket
	result := s.DB.
		Clauses(clause.Returning{}).
		Delete(&market, "security_uuid = ? AND market_code = ?", securityUuid, marketCode)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		return nil, model.ErrNotFound
	}
	return s.securityMarketModelFromDb(market), nil
}

// UpdateSecurityTaxonomies creates/updates/deletes taxonomies of security
func (s *securityService) UpdateSecurityTaxonomies(
	securityUuid, rootTaxonomyUuid uuid.UUID, inputs []*model.SecurityTaxonomyInput,
) (
	[]*model.SecurityTaxonomy, error,
) {
	// Remove securityTaxonomies of rootTaxonomy not in inputs
	secTaxonomyUuids := make([]uuid.UUID, len(inputs))
	for i := range inputs {
		secTaxonomyUuids[i] = inputs[i].TaxonomyUUID
	}
	var err error
	if len(secTaxonomyUuids) == 0 {
		err = s.DB.Exec("DELETE FROM securities_taxonomies st "+
			"USING taxonomies t "+
			"WHERE st.taxonomy_uuid = t.uuid"+
			" AND st.security_uuid = ?"+
			" AND t.root_uuid = ?", securityUuid, rootTaxonomyUuid).
			Error
	} else {
		err = s.DB.Exec("DELETE FROM securities_taxonomies st "+
			"USING taxonomies t "+
			"WHERE st.taxonomy_uuid = t.uuid"+
			" AND st.security_uuid = ?"+
			" AND t.root_uuid = ?"+
			" AND st.taxonomy_uuid NOT IN ?", securityUuid, rootTaxonomyUuid, secTaxonomyUuids).
			Error
	}
	if err != nil {
		panic(err)
	}

	// Upsert all security taxonomies in input
	upsert := make([]db.SecurityTaxonomy, len(inputs))
	for i := range inputs {
		upsert[i].SecurityUUID = securityUuid
		upsert[i].TaxonomyUUID = inputs[i].TaxonomyUUID
		upsert[i].Weight = inputs[i].Weight
	}
	if len(upsert) > 0 {
		if err := s.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&upsert).Error; err != nil {
			panic(err)
		}
	}

	return s.securityTaxonomiesModelFromDb(upsert), nil
}

// UpsertTag creates/updates tag
func (s *securityService) UpsertTag(name string, securityUuids []uuid.UUID) ([]*model.Security, error) {

	// Get or create tag
	var tag db.Tag
	if err := s.DB.
		Where("LOWER(name) = LOWER(?)", name).
		Attrs(db.Tag{UUID: uuid.New()}).
		FirstOrCreate(&tag, db.Tag{Name: name}).Error; err != nil {
		panic(err)
	}

	// Delete removed associations
	if err := s.DB.Delete(&db.SecurityTag{}, "tag_uuid = ? AND security_uuid NOT IN ?", tag.UUID, securityUuids).Error; err != nil {
		panic(err)
	}

	// Create new associations
	upsert := make([]db.SecurityTag, len(securityUuids))
	for i := range securityUuids {
		upsert[i].TagUUID = tag.UUID
		upsert[i].SecurityUUID = securityUuids[i]
	}
	if len(upsert) > 0 {
		if err := s.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&upsert).Error; err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
				return nil, fmt.Errorf("data violates constraint " + pqErr.Constraint)
			}

			panic(err)
		}
	}

	// Get associated securities
	var dbSecurities []db.Security
	if err := s.DB.Model(&tag).Association("Securities").Find(&dbSecurities); err != nil {
		panic(err)
	}

	securities := make([]*model.Security, len(dbSecurities))
	for i := range dbSecurities {
		securities[i] = s.modelFromDb(dbSecurities[i])
	}

	return securities, nil
}

// DeleteTag removes tag
func (s *securityService) DeleteTag(name string) {
	if err := s.DB.Where("LOWER(name) = LOWER(?)", name).Delete(&db.Tag{}).Error; err != nil {
		panic(err)
	}
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

// securityTaxonomiesModelFromDb converts list of security taxonomies from database into model
func (*securityService) securityTaxonomiesModelFromDb(secTaxonomies []db.SecurityTaxonomy) []*model.SecurityTaxonomy {
	ret := make([]*model.SecurityTaxonomy, len(secTaxonomies))
	for i := range secTaxonomies {
		ret[i] = &model.SecurityTaxonomy{
			SecurityUUID: secTaxonomies[i].SecurityUUID,
			TaxonomyUUID: secTaxonomies[i].TaxonomyUUID,
			Weight:       secTaxonomies[i].Weight,
		}
	}
	return ret
}

// securityMarketModelFromDb converts security market from database into model
func (*securityService) securityMarketModelFromDb(m db.SecurityMarket) *model.SecurityMarket {
	return &model.SecurityMarket{
		MarketCode:     m.MarketCode,
		CurrencyCode:   m.CurrencyCode,
		Symbol:         m.Symbol,
		FirstPriceDate: m.FirstPriceDate,
		LastPriceDate:  m.LastPriceDate,
	}
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
