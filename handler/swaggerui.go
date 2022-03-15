package handler

import (
	_ "embed"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/swaggerui"
)

//go:embed openapi.json
var openApiSpec []byte

func (h *Handler) RegisterSwaggerUi(g *gin.RouterGroup, prefix string) {
	ginHandler := gin.WrapH(
		http.StripPrefix(path.Join(g.BasePath(), prefix), swaggerui.Handler(openApiSpec)),
	)
	g.GET(prefix+"/*w", ginHandler)
}
