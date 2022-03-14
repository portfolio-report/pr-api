package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/models"
)

type contextKey struct {
	name string
}

// Private key for context to prevent possible collisions
var userCtxKey = &contextKey{name: "user"}

// Reads authorization token from HTTP header (if any),
// stores user in request context if token corresponds to valid session
func AuthUser(s models.SessionService, u models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := s.GetSessionToken(c)

		if token == "" {
			c.Next()
			return
		}

		session, err := s.ValidateToken(token)
		if err != nil {
			panic(err)
		}

		if session == nil {
			c.Next()
			return
		}

		user, err := u.GetUserFromSession(session)
		if err != nil {
			panic(err)
		}

		go func() {
			err := u.UpdateLastSeen(user)
			if err != nil {
				fmt.Println("Error in background processing:", err.Error())
			}
		}()

		ctx := context.WithValue(c.Request.Context(), userCtxKey, user)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// Gets logged in user from request context,
// may contain nil value, if no user is logged in.
func UserFromContext(ctx context.Context) *model.User {
	v := ctx.Value(userCtxKey)
	user, valid := v.(*model.User)
	if valid {
		return user
	} else {
		return nil
	}
}
