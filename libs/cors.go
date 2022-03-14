package libs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Adds CORS header to all responses and handles CORS preflight checks
//
// Behaviour is independent of request headers (Origin or Access-Control-*) to allow caching
func Cors(c *gin.Context) {
	headers := c.Writer.Header()

	headers.Set("Access-Control-Allow-Origin", "*")

	// Handle preflight check
	if c.Request.Method == "OPTIONS" {
		headers.Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
		headers.Set("Access-Control-Allow-Headers", "Origin,Content-Length,Content-Type,Authorization")
		headers.Set("Access-Control-Max-Age", "86400") // 1d
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
}
