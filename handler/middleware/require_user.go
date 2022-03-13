package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
	"github.com/portfolio-report/pr-api/models"
)

// returns HTTP error if no user is logged in
//
// requires middleware AuthUser to be run
func RequireUser(s models.SessionService, u models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserFromContext(c.Request.Context())

		if user == nil {
			libs.HandleUnauthorizedError(c)
			return
		}

		c.Next()
	}
}
