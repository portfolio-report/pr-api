package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs"
)

// RequireUser returns middleware which only passes if a user is logged in
//
// Requires middleware AuthUser to be run
func RequireUser(s model.SessionService, u model.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserFromContext(c.Request.Context())

		if user == nil {
			libs.HandleUnauthorizedError(c)
			return
		}

		c.Next()
	}
}
