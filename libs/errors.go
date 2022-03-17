package libs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleForbiddenError returns HTTP Forbidden error with JSON body
func HandleForbiddenError(c *gin.Context, msg string) {
	code := http.StatusForbidden
	c.JSON(code, gin.H{"statusCode": code, "error": http.StatusText(code), "message": msg})
	c.Abort()
}

// HandleUnauthorizedError returns Unauthorized error with JSON body
func HandleUnauthorizedError(c *gin.Context) {
	code := http.StatusUnauthorized
	c.JSON(code, gin.H{"statusCode": code, "error": http.StatusText(code)})
	c.Abort()
}

// HandleNotFoundError returns Not Found error with JSON body
func HandleNotFoundError(c *gin.Context) {
	code := http.StatusNotFound
	c.JSON(code, gin.H{"statusCode": code, "error": http.StatusText(code)})
	c.Abort()
}

// HandleBadRequestError returns Bad Request error with JSON body
func HandleBadRequestError(c *gin.Context, msg string) {
	code := http.StatusBadRequest
	c.JSON(code, gin.H{"statusCode": code, "error": http.StatusText(code), "message": msg})
	c.Abort()
}
