package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type SecurityServiceTestSuite struct {
	suite.Suite
	db      *gorm.DB
	service *securityService
}

func (s *SecurityServiceTestSuite) SetupSuite() {
	godotenv.Load("../.env")

	var err error
	s.db, err = db.InitDb(ReadConfig().Db)
	s.Nil(err)

	service := NewSecurityService(s.db)
	var ok bool
	s.service, ok = service.(*securityService)
	s.True(ok)
}

func (s *SecurityServiceTestSuite) TearDownSuite() {
	sql, err := s.db.DB()
	s.Nil(err)
	sql.Close()
}

func TestSecurityService(t *testing.T) {
	suite.Run(t, new(SecurityServiceTestSuite))
}

func (s *SecurityServiceTestSuite) TestGetSecurityByUUID() {
	dbSecurity := db.Security{UUID: uuid.New().String()}
	err := s.db.Create(&dbSecurity).Error
	s.Nil(err)

	{
		security, err := s.service.GetSecurityByUUID(dbSecurity.UUID)
		s.Nil(err)
		s.NotNil(security)
	}

	{
		_, err := s.service.GetSecurityByUUID("952df501-1e22-4693-a208-0c013cb1b415")
		s.ErrorIs(err, gorm.ErrRecordNotFound)
	}
}

func (s *SecurityServiceTestSuite) TestGetEventsOfSecurity() {
	dbSecurity := db.Security{UUID: uuid.New().String()}
	err := s.db.Create(&dbSecurity).Error
	s.Nil(err)
	security := model.Security{UUID: dbSecurity.UUID}

	{
		events, err := s.service.GetEventsOfSecurity(&security)
		s.Nil(err)
		s.Len(events, 0)
	}
}
