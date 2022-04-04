package service

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type taxonomyService struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

// NewTaxonomyService creates and returns new taxonomy service
func NewTaxonomyService(db *gorm.DB, validate *validator.Validate) model.TaxonomyService {
	return &taxonomyService{
		DB:       db,
		Validate: validate,
	}
}

// modelFromDb converts taxonomy from database to model
func (*taxonomyService) modelFromDb(t db.Taxonomy) *model.Taxonomy {
	return &model.Taxonomy{
		UUID:       t.UUID,
		ParentUUID: t.ParentUUID,
		RootUUID:   t.RootUUID,
		Name:       t.Name,
		Code:       t.Code,
	}
}

// GetAllTaxonomies returns all taxonomies
func (s *taxonomyService) GetAllTaxonomies() []*model.Taxonomy {
	var taxonomies []db.Taxonomy
	if err := s.DB.Find(&taxonomies).Error; err != nil {
		panic(err)
	}

	ret := []*model.Taxonomy{}
	for _, t := range taxonomies {
		ret = append(ret, s.modelFromDb(t))
	}

	return ret
}

// GetTaxonomyByUUID returns taxonomy identified by UUID
func (s *taxonomyService) GetTaxonomyByUUID(uuid uuid.UUID) (*model.Taxonomy, error) {
	var t db.Taxonomy
	err := s.DB.Take(&t, "uuid = ?", uuid).Error
	return s.modelFromDb(t), err
}

// GetDescendantsOfTaxonomy returns all descendants of taxonomy,
// i.e. children, children of children, etc.
func (s *taxonomyService) GetDescendantsOfTaxonomy(taxonomy *model.Taxonomy) []*model.Taxonomy {
	var taxonomies []db.Taxonomy
	err := s.DB.Find(&taxonomies, "root_uuid = ?", taxonomy.UUID).Error
	if err != nil {
		panic(err)
	}

	ret := []*model.Taxonomy{}
	for _, t := range taxonomies {
		ret = append(ret, s.modelFromDb(t))
	}

	return ret
}

// CreateTaxonomy creates new taxonomy
func (s *taxonomyService) CreateTaxonomy(input *model.TaxonomyInput) (*model.Taxonomy, error) {
	taxonomy := db.Taxonomy{
		UUID:       uuid.New(),
		Name:       input.Name,
		Code:       input.Code,
		ParentUUID: input.ParentUUID,
	}

	if input.ParentUUID != nil {
		parent, err := s.GetTaxonomyByUUID(*input.ParentUUID)
		if err != nil {
			return nil, fmt.Errorf("parentUuid invalid")
		}

		if parent.RootUUID != nil {
			taxonomy.RootUUID = parent.RootUUID
		} else {
			taxonomy.RootUUID = input.ParentUUID
		}
	}

	if err := s.DB.Clauses(clause.Returning{}).Create(&taxonomy).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates constraint " + pqErr.Constraint)
		}

		panic(err)
	}

	return s.modelFromDb(taxonomy), nil
}

// UpdateTaxonomy updates taxonomy
func (s *taxonomyService) UpdateTaxonomy(uuid uuid.UUID, input *model.TaxonomyInput) (*model.Taxonomy, error) {
	var taxonomy db.Taxonomy
	if err := s.DB.Take(&taxonomy, "uuid = ?", uuid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		panic(err)
	}

	taxonomy.Name = input.Name
	taxonomy.Code = input.Code
	taxonomy.ParentUUID = input.ParentUUID

	if input.ParentUUID != nil {
		if *input.ParentUUID == uuid {
			return nil, fmt.Errorf("parentUuid must be different from own uuid")
		}

		parent, err := s.GetTaxonomyByUUID(*input.ParentUUID)
		if err != nil {
			return nil, fmt.Errorf("parentUuid invalid")
		}

		if parent.RootUUID != nil {
			taxonomy.RootUUID = parent.RootUUID
		} else {
			taxonomy.RootUUID = input.ParentUUID
		}
	} else {
		taxonomy.RootUUID = nil
	}

	if err := s.DB.Clauses(clause.Returning{}).Save(&taxonomy).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates constraint " + pqErr.Constraint)
		}

		panic(err)
	}

	return s.modelFromDb(taxonomy), nil
}

// DeleteTaxonomy deletes taxonomy
func (s *taxonomyService) DeleteTaxonomy(uuid uuid.UUID) (*model.Taxonomy, error) {
	taxonomy := db.Taxonomy{}
	result := s.DB.Clauses(clause.Returning{}).Delete(&taxonomy, "uuid = ?", uuid)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("Not found")
	}

	return s.modelFromDb(taxonomy), nil
}
