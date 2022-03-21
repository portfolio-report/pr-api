package service

import (
	"context"
	"testing"
	"time"

	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type UserServiceTestSuite struct {
	suite.Suite
	db      *gorm.DB
	service *userService
}

func (s *UserServiceTestSuite) SetupSuite() {
	godotenv.Load("../.env")

	var err error
	s.db, err = db.InitDb(ReadConfig().Db)
	s.Nil(err)

	service := NewUserService(s.db)
	var ok bool
	s.service, ok = service.(*userService)
	s.True(ok)
}

func (s *UserServiceTestSuite) TearDownSuite() {
	s.db.Delete(&db.User{}, "username = 'testuser'")

	sql, err := s.db.DB()
	s.Nil(err)
	sql.Close()
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) TestUserLifecycle() {
	// Create user
	user, err := s.service.Create("TestUser")
	s.Nil(err)
	s.IsType(&model.User{}, user)
	s.Equal("testuser", user.Username)
	s.False(user.IsAdmin)
	s.Greater(user.ID, 0)

	var dbUser db.User
	err = s.db.Take(&dbUser, "username = 'testuser'").Error
	s.Nil(err)
	s.Nil(dbUser.Password)

	// Create already existing user
	{
		user, err := s.service.Create("testuser")
		s.Nil(user)
		s.ErrorIs(err, ErrUserExistsAlready)
	}

	// Get existing user
	{
		user, err := s.service.GetUserByUsername(context.TODO(), "tesTuseR")
		s.Nil(err)
		s.IsType(&model.User{}, user)
		s.Equal("testuser", user.Username)
		s.False(user.IsAdmin)
		s.Greater(user.ID, 0)
	}

	// Get non-existent user
	{
		_, err := s.service.GetUserByUsername(context.TODO(), "unknown-user")
		s.ErrorIs(err, gorm.ErrRecordNotFound)
	}

	// Verify non-existent password
	{
		ok, err := s.service.VerifyPassword(context.TODO(), user, "testpassword")
		s.Nil(err)
		s.False(ok)
	}

	// Update password
	{
		err := s.service.UpdatePassword(context.TODO(), user, "testpassword")
		s.Nil(err)
		err = s.db.Take(&dbUser, "username = 'testuser'").Error
		s.Nil(err)
		s.NotNil(dbUser.Password)
	}

	// Verify password
	{
		ok, err := s.service.VerifyPassword(context.TODO(), user, "testpassword")
		s.Nil(err)
		s.True(ok)
	}
	{
		ok, err := s.service.VerifyPassword(context.TODO(), user, "wrong")
		s.Nil(err)
		s.False(ok)
	}

	// Update last seen
	{
		err := s.service.UpdateLastSeen(user)
		s.Nil(err)
		err = s.db.Take(&dbUser, "username = 'testuser'").Error
		s.Nil(err)
		s.True(time.Time(dbUser.LastSeenAt).After(time.Now().AddDate(0, 0, -1)))
	}

	// Delete user
	{
		err = s.service.Delete(user.ID)
		s.Nil(err)

		err = s.db.Take(&dbUser, "username = 'testuser'").Error
		s.ErrorIs(err, gorm.ErrRecordNotFound)
	}
}
