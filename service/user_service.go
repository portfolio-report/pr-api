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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userService struct {
	DB *gorm.DB
}

// NewUserService creates and returns new user service
func NewUserService(db *gorm.DB) model.UserService {
	return &userService{
		DB: db,
	}
}

// ErrUserExistsAlready indicates a user could not be created,
// because the username is used already
var ErrUserExistsAlready = errors.New("user exists already")

// modelFromDb converts user from database into model
func (*userService) modelFromDb(u db.User) *model.User {
	return &model.User{
		ID:         int(u.ID),
		Username:   u.Username,
		IsAdmin:    u.IsAdmin,
		LastSeenAt: time.Time(u.LastSeenAt).Format("2006-01-02"),
	}
}

// Create creates user (without password)
func (s *userService) Create(username string) (*model.User, error) {
	user := db.User{
		Username: strings.ToLower(username),
	}

	if err := s.DB.Clauses(clause.Returning{}).Create(&user).Error; err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, ErrUserExistsAlready
		}

		panic(err)
	}

	return s.modelFromDb(user), nil
}

// GetUserByUsername return user identified by username
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	username = strings.ToLower(username)
	var user db.User
	err := s.DB.Take(&user, "username = ?", username).Error
	return s.modelFromDb(user), err
}

// GetByIDs returns users by ids
func (s *userService) GetByIDs(ids []int) ([]*model.User, error) {
	var users []db.User
	err := s.DB.Find(&users, "id IN ?", ids).Error
	if err != nil {
		panic(err)
	}

	// map to model
	result := make([]*model.User, len(users))
	for i := range users {
		result[i] = s.modelFromDb(users[i])
	}
	return result, nil
}

// GetUserFromSession returns user that owns session
func (s *userService) GetUserFromSession(session *model.Session) (*model.User, error) {
	var user db.User
	err := s.DB.Take(&user, session.UserID).Error
	if err != nil {
		panic(err)
	}
	return s.modelFromDb(user), nil
}

// UpdatePassword changes the user's password
func (s *userService) UpdatePassword(ctx context.Context, user *model.User, password string) error {
	hash, err := argon2.HashPasswordDefault(password)
	if err != nil {
		panic(err)
	}

	if err := s.DB.Model(db.User{ID: uint(user.ID)}).Update("password", hash).Error; err != nil {
		panic(err)
	}

	return nil
}

// VerifyPassword checks if password matches the hash stored in database
func (s *userService) VerifyPassword(ctx context.Context, user *model.User, password string) (bool, error) {
	var dbUser db.User
	err := s.DB.Take(&dbUser, user.ID).Error
	if err != nil {
		panic(err)
	}
	if dbUser.Password == nil {
		return false, nil
	}
	return argon2.VerifyPassword(password, *dbUser.Password)
}

// Delete removes user
func (s *userService) Delete(id int) error {
	return s.DB.Delete(db.User{}, "id = ?", id).Error
}

// UpdateLastSeen updates the date a user was seen last
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
