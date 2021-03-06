package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Private key for context to prevent possible collisions
var useragentCtxKey = &contextKey{name: "useragent"}

// Useragent stores user agent in request context
func Useragent(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), useragentCtxKey, c.GetHeader("User-Agent"))
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

// UseragentFromContext retrieves user agent from request context,
// panics if context entry does not exists (middleware not run).
func UseragentFromContext(ctx context.Context) string {
	return ctx.Value(useragentCtxKey).(string)
}
