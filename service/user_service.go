package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs/argon2"
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) models.UserService {
	return &userService{
		DB: db,
	}
}

var UserExistsAlreadyError = errors.New("user exists already")

func (*userService) ModelFromDb(u db.User) *model.User {
	return &model.User{
		ID:         int(u.ID),
		Username:   u.Username,
		IsAdmin:    u.IsAdmin,
		LastSeenAt: time.Time(u.LastSeenAt).Format("2006-01-02"),
	}
}

func (s *userService) Create(username string) (*model.User, error) {
	user := db.User{
		Username: strings.ToLower(username),
	}

	if err := s.DB.Clauses(clause.Returning{}).Create(&user).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, UserExistsAlreadyError
		}

		panic(err)
	}

	return s.ModelFromDb(user), nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	username = strings.ToLower(username)
	var user db.User
	err := s.DB.Take(&user, "username = ?", username).Error
	return s.ModelFromDb(user), err
}

func (s *userService) GetUserFromSession(session *model.Session) *model.User {
	var user db.User
	err := s.DB.Take(&user, session.UserID).Error
	if err != nil {
		panic(err)
	}
	return s.ModelFromDb(user)
}

func (s *userService) UpdatePassword(ctx context.Context, user *model.User, password string) error {
	hash, err := argon2.HashPasswordDefault(password)
	if err != nil {
		panic(err)
	}

	if err := s.DB.Model(db.User{ID: uint(user.ID)}).Update("password", hash).Error; err != nil {
		return err
	}

	return nil
}

func (s *userService) VerifyPassword(ctx context.Context, user *model.User, password string) (bool, error) {
	var dbUser db.User
	err := s.DB.Take(&dbUser, user.ID).Error
	if err != nil {
		panic(err)
	}
	return argon2.VerifyPassword(password, *dbUser.Password)
}

func (s *userService) Delete(user *model.User) error {
	return s.DB.Delete(db.User{}, "id = ?", user.ID).Error
}

func (s *userService) UpdateLastSeen(user *model.User) error {
	now := time.Now().Format("2006-01-02")

	if now != user.LastSeenAt {
		return s.DB.
			Table("users").
			Where("id = ?", user.ID).
			Update("last_seen_at", now).Error
	}
	return nil
}
