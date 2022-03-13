package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type taxonomyService struct {
	DB *gorm.DB
}

func NewTaxonomyService(db *gorm.DB) models.TaxonomyService {
	return &taxonomyService{
		DB: db,
	}
}

func (*taxonomyService) ModelFromDb(t db.Taxonomy) *model.Taxonomy {
	return &model.Taxonomy{
		UUID:       t.UUID,
		ParentUUID: t.ParentUUID,
		RootUUID:   t.RootUUID,
		Name:       t.Name,
		Code:       t.Code,
	}
}

func (s *taxonomyService) GetAllTaxonomies() ([]*model.Taxonomy, error) {
	var taxonomies []db.Taxonomy
	if err := s.DB.Find(&taxonomies).Error; err != nil {
		panic(err)
	}

	ret := []*model.Taxonomy{}
	for _, t := range taxonomies {
		ret = append(ret, s.ModelFromDb(t))
	}

	return ret, nil
}

func (s *taxonomyService) GetTaxonomyByUUID(uuid string) (*model.Taxonomy, error) {
	var t db.Taxonomy
	err := s.DB.Take(&t, "uuid = ?", uuid).Error
	return s.ModelFromDb(t), err
}

func (s *taxonomyService) GetDescendantsOfTaxonomy(taxonomy *model.Taxonomy) ([]*model.Taxonomy, error) {
	var taxonomies []db.Taxonomy
	err := s.DB.Find(&taxonomies, "root_uuid = ?", taxonomy.UUID).Error
	if err != nil {
		panic(err)
	}

	ret := []*model.Taxonomy{}
	for _, t := range taxonomies {
		ret = append(ret, s.ModelFromDb(t))
	}

	return ret, nil
}

func (s *taxonomyService) CreateTaxonomy(input *model.Taxonomy) (*model.Taxonomy, error) {
	taxonomy := db.Taxonomy{
		UUID: uuid.New().String(),
		Name: input.Name,
		Code: input.Code,
	}

	if input.ParentUUID != nil {
		parent, err := s.GetTaxonomyByUUID(*input.ParentUUID)
		if err != nil {
			return nil, fmt.Errorf("parentUuid invalid")
		}

		taxonomy.ParentUUID = input.ParentUUID

		if parent.RootUUID != nil {
			taxonomy.RootUUID = parent.RootUUID
		} else {
			taxonomy.RootUUID = input.ParentUUID
		}
	}

	if err := s.DB.Create(&taxonomy).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates constraint " + pqErr.Constraint)
		}

		panic(err)
	}

	return s.ModelFromDb(taxonomy), nil
}

func (s *taxonomyService) UpdateTaxonomy(uuid string, input *model.Taxonomy) (*model.Taxonomy, error) {

	var taxonomy db.Taxonomy
	if err := s.DB.Take(&taxonomy, "uuid = ?", uuid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		panic(err)
	}

	if input.Name != "" {
		taxonomy.Name = input.Name
	}

	if input.ParentUUID != nil {
		if *input.ParentUUID == uuid {
			return nil, fmt.Errorf("parentUuid must be different from own uuid")
		}

		if *input.ParentUUID == "" {
			taxonomy.ParentUUID = nil
			taxonomy.RootUUID = nil
		} else {
			parent, err := s.GetTaxonomyByUUID(*input.ParentUUID)
			if err != nil {
				return nil, fmt.Errorf("parentUuid invalid")
			}

			taxonomy.ParentUUID = input.ParentUUID

			if parent.RootUUID != nil {
				taxonomy.RootUUID = parent.RootUUID
			} else {
				taxonomy.RootUUID = input.ParentUUID
			}
		}
	}

	if input.Code != nil {
		if *input.Code == "" && taxonomy.ParentUUID == nil {
			// Root nodes may have empty code
			input.Code = nil
		} else {
			// Non-root nodes must have an unique code
			taxonomy.Code = input.Code
		}
	}

	if err := s.DB.Clauses(clause.Returning{}).Save(&taxonomy).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return nil, fmt.Errorf("data violates constraint " + pqErr.Constraint)
		}

		panic(err)
	}

	return s.ModelFromDb(taxonomy), nil
}

func (s *taxonomyService) DeleteTaxonomy(uuid string) (*model.Taxonomy, error) {
	taxonomy := db.Taxonomy{}
	result := s.DB.Clauses(clause.Returning{}).Delete(&taxonomy, "uuid = ?", uuid)
	if err := result.Error; err != nil {
		panic(err)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("Not found")
	}

	return s.ModelFromDb(taxonomy), nil
}
