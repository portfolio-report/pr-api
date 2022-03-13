package libs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleForbiddenError(c *gin.Context, msg string) {
	code := http.StatusForbidden
	c.JSON(code, gin.H{"statusCode": code, "error": http.StatusText(code), "message": msg})
	c.Abort()
}

func HandleUnauthorizedError(c *gin.Context) {
	code := http.StatusUnauthorized
	c.JSON(code, gin.H{"statusCode": code, "error": http.StatusText(code)})
	c.Abort()
}

func HandleNotFoundError(c *gin.Context) {
	code := http.StatusNotFound
	c.JSON(code, gin.H{"statusCode": code, "error": http.StatusText(code)})
	c.Abort()
}

func HandleBadRequestError(c *gin.Context, msg string) {
	code := http.StatusBadRequest
	c.JSON(code, gin.H{"statusCode": code, "error": http.StatusText(code), "message": msg})
	c.Abort()
}
