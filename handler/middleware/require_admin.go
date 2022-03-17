package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/portfolio-report/pr-api/libs"
)

// RequireAdmin returns middleware which only passes if user has admin privileges
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserFromContext(c.Request.Context())
		if !user.IsAdmin {
			libs.HandleUnauthorizedError(c)
			return
		}

		c.Next()
	}
}
