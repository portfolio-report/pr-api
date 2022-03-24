package handler

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/libs"
)

// GetRoot returns static ok message
func (h *rootHandler) GetRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "message": "ok"})
}

// GetVersion returns version information from build
func (h *rootHandler) GetVersion(c *gin.Context) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		libs.HandleNotFoundError(c)
		return
	}

	ret := gin.H{}

	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			ret["revision"] = kv.Value
		case "vcs.time":
			ret["time"] = kv.Value
		case "vcs.modified":
			ret["modified"] = kv.Value
		}
	}
	c.JSON(http.StatusOK, ret)
}
