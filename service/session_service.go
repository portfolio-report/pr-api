package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/portfolio-report/pr-api/db"
	"github.com/portfolio-report/pr-api/graph/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sessionService struct {
	DB       *gorm.DB
	Validate *validator.Validate
	Timeout  time.Duration
}

// NewSessionService creates and returns new session service
func NewSessionService(db *gorm.DB, validate *validator.Validate, timeout time.Duration) model.SessionService {
	return &sessionService{
		DB:       db,
		Validate: validate,
		Timeout:  timeout,
	}
}

// modelFromDb converts session from database into model
func (*sessionService) modelFromDb(s db.Session) *model.Session {
	return &model.Session{
		Token:          s.Token,
		Note:           s.Note,
		CreatedAt:      s.CreatedAt.UTC(),
		LastActivityAt: s.LastActivityAt.UTC(),
		UserID:         s.UserID,
	}
}

// GetAllOfUser returns all sessions of user
func (s *sessionService) GetAllOfUser(user *model.User) []*model.Session {
	var sessions []db.Session
	err := s.DB.Find(&sessions, "user_id = ?", user.ID).Error
	if err != nil {
		panic(err)
	}

	response := []*model.Session{}
	for _, session := range sessions {
		response = append(response, s.modelFromDb(session))
	}
	return response
}

// CreateSession creates new session for user
func (s *sessionService) CreateSession(user *model.User, note string) (*model.Session, error) {
	token := uuid.New().String()

	session := db.Session{
		Token:  token,
		Note:   note,
		UserID: uint(user.ID),
	}
	err := s.DB.
		Select("Token", "Note", "UserID"). // only insert certain columns
		Clauses(clause.Returning{}).       // return db defaults for remaining columns
		Create(&session).
		Error

	return s.modelFromDb(session), err
}

// DeleteSession removes session
func (s *sessionService) DeleteSession(token string) (*model.Session, error) {
	var session db.Session
	err := s.DB.Clauses(clause.Returning{}).Where("token = ?", token).Delete(&session).Error
	return s.modelFromDb(session), err
}

// sessionLastActivityLimit returns the limit for last activitiy of a session
// to be still considered valid
func (s *sessionService) sessionLastActivityLimit() time.Time {
	return time.Now().Add(-1 * s.Timeout)
}

type authHeader struct {
	IDToken string `header:"Authorization"`
}

// GetSessionToken returns session token from HTTP headers
func (*sessionService) GetSessionToken(c *gin.Context) string {
	h := authHeader{}

	if err := c.ShouldBindHeader(&h); err != nil {
		return ""
	}

	idTokenHeader := strings.Split(h.IDToken, "Bearer ")

	if len(idTokenHeader) < 2 {
		return ""
	}

	return idTokenHeader[1]
}

// ValidateToken checks if session token is valid.
// Returns session for of valid tokens, nil for invalid session tokens.
func (s *sessionService) ValidateToken(token string) (*model.Session, error) {
	var session db.Session

	if err := s.Validate.Var(token, "uuid"); err != nil {
		return nil, nil
	}

	err := s.DB.
		Where("last_activity_at > ?", s.sessionLastActivityLimit()).
		Where("token = ?", token).
		Take(&session).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		panic(err)
	}

	go func() {
		err := s.updateLastActivity(&session)
		if err != nil {
			fmt.Println("Error in background processing:", err.Error())
		}
	}()

	return s.modelFromDb(session), nil
}

// updateLastActivity sets last activity of session to now
func (s *sessionService) updateLastActivity(session *db.Session) error {
	now := time.Now()
	if now.Sub(session.LastActivityAt).Seconds() > 60 {
		return s.DB.
			Table("sessions").
			Where("token = ?", session.Token).
			Update("last_activity_at", now).Error
	}
	return nil
}

// CleanupExpiredSessions removes expired sessions from database
func (s *sessionService) CleanupExpiredSessions() error {
	result := s.DB.
		Where("last_activity_at < ?", s.sessionLastActivityLimit()).
		Delete(&db.Session{})
	log.Printf("removed %d expired sessions\n", result.RowsAffected)
	return result.Error
}
