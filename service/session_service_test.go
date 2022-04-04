package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type SessionServiceTestSuite struct {
	suite.Suite
	db      *gorm.DB
	service *sessionService
	dbUser  *db.User
	user    *model.User
}

func (s *SessionServiceTestSuite) SetupSuite() {
	godotenv.Load("../.env")

	var err error
	s.db, err = db.InitDb(ReadConfig().Db)
	s.Nil(err)

	service := NewSessionService(s.db, validator.New(), 5*time.Second)
	var ok bool
	s.service, ok = service.(*sessionService)
	s.True(ok)

	s.dbUser = &db.User{Username: "testuser-session"}
	err = s.db.Create(s.dbUser).Error
	s.Nil(err)
	s.user = &model.User{ID: int(s.dbUser.ID), Username: s.dbUser.Username}
}

func (s *SessionServiceTestSuite) TearDownSuite() {
	s.db.Delete(&db.User{}, "username = 'testuser-session'")

	sql, err := s.db.DB()
	s.Nil(err)
	sql.Close()
}

func TestSessionService(t *testing.T) {
	suite.Run(t, new(SessionServiceTestSuite))
}

func (s *SessionServiceTestSuite) TestGetAllOfUser() {
	err := s.db.Delete(&db.Session{}, "user_id = ?", s.dbUser.ID).Error
	s.Nil(err)

	sessions := s.service.GetAllOfUser(s.user)
	s.Len(sessions, 0)

	session, err := s.service.CreateSession(s.user, "some note")
	s.Nil(err)

	sessions = s.service.GetAllOfUser(s.user)
	s.Len(sessions, 1)
	s.Equal(session, sessions[0])
}

func (s *SessionServiceTestSuite) TestCreateSession() {
	session, err := s.service.CreateSession(s.user, "some note")
	s.Nil(err)
	s.IsType(&model.Session{}, session)
	s.Equal(uint(s.user.ID), session.UserID)
	s.Equal("some note", session.Note)
	s.True(time.Time(session.CreatedAt).After(time.Now().Add(-1 * time.Second)))
	s.True(time.Time(session.CreatedAt).Before(time.Now().Add(1 * time.Second)))
	s.True(time.Time(session.LastActivityAt).After(time.Now().Add(-1 * time.Second)))
	s.True(time.Time(session.LastActivityAt).Before(time.Now().Add(1 * time.Second)))

	err = s.db.Take(&db.Session{}, "token = ?", session.Token).Error
	s.Nil(err)
}

func (s *SessionServiceTestSuite) TestDeleteSession() {
	session, err := s.service.CreateSession(s.user, "some note")
	s.Nil(err)

	_, err = s.service.DeleteSession(session.Token)
	s.Nil(err)

	err = s.db.Take(&db.Session{}, "token = ?", session.Token).Error
	s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *SessionServiceTestSuite) TestValidateToken() {
	origSession, err := s.service.CreateSession(s.user, "some note")
	s.Nil(err)

	session, err := s.service.ValidateToken(origSession.Token)
	s.Nil(err)
	s.Equal(origSession, session)

	session, err = s.service.ValidateToken("foobar")
	s.Nil(err)
	s.Nil(session)

	session, err = s.service.ValidateToken("6389ddeb-b9a2-4fc9-95f9-8cd6f0505b9b")
	s.Nil(err)
	s.Nil(session)
}

func (s *SessionServiceTestSuite) TestGetSessionToken() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	{
		c.Request = &http.Request{
			Header: http.Header{
				"Authorization": []string{"Bearer foo"},
			},
		}
		token := s.service.GetSessionToken(c)
		s.Equal("foo", token)
	}

	{
		c.Request = &http.Request{
			Header: http.Header{
				"Authorization": []string{"messed up"},
			},
		}
		token := s.service.GetSessionToken(c)
		s.Equal("", token)
	}

	{
		c.Request = &http.Request{
			Header: http.Header{},
		}
		token := s.service.GetSessionToken(c)
		s.Equal("", token)
	}
}
