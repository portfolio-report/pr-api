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
	"github.com/portfolio-report/pr-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sessionService struct {
	DB       *gorm.DB
	Validate *validator.Validate
	Timeout  int
}

func NewSessionService(db *gorm.DB, validate *validator.Validate, timeout int) models.SessionService {
	return &sessionService{
		DB:       db,
		Validate: validate,
		Timeout:  timeout,
	}
}

func (*sessionService) ModelFromDb(s db.Session) *model.Session {
	return &model.Session{
		Token:          s.Token,
		Note:           s.Note,
		CreatedAt:      s.CreatedAt.UTC(),
		LastActivityAt: s.LastActivityAt.UTC(),
		UserID:         s.UserID,
	}
}

func (s *sessionService) GetAllOfUser(user *model.User) []*model.Session {
	var sessions []db.Session
	err := s.DB.Find(&sessions, "user_id = ?", user.ID).Error
	if err != nil {
		panic(err)
	}

	response := []*model.Session{}
	for _, session := range sessions {
		response = append(response, s.ModelFromDb(session))
	}
	return response
}

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

	return s.ModelFromDb(session), err
}

func (s *sessionService) DeleteSession(token string) (*model.Session, error) {
	var session db.Session
	err := s.DB.Clauses(clause.Returning{}).Where("token = ?", token).Delete(&session).Error
	return s.ModelFromDb(session), err
}

func (s *sessionService) sessionLastActivityLimit() time.Time {
	return time.Now().Add(-1 * time.Duration(s.Timeout) * time.Second)
}

type authHeader struct {
	IDToken string `header:"Authorization"`
}

func (s *sessionService) GetSessionToken(c *gin.Context) string {
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

func (s *sessionService) ValidateToken(token string) *model.Session {
	var session db.Session

	if err := s.Validate.Var(token, "uuid"); err != nil {
		return nil
	}

	err := s.DB.
		Where("last_activity_at > ?", s.sessionLastActivityLimit()).
		Where("token = ?", token).
		Take(&session).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		panic(err)
	}

	go func() {
		err := s.updateLastActivity(&session)
		if err != nil {
			fmt.Println("Error in background processing:", err.Error())
		}
	}()

	return s.ModelFromDb(session)
}

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

func (s *sessionService) CleanupExpiredSessions() error {
	result := s.DB.
		Where("last_activity_at < ?", s.sessionLastActivityLimit()).
		Delete(&db.Session{})
	log.Printf("removed %d expired sessions\n", result.RowsAffected)
	return result.Error
}
